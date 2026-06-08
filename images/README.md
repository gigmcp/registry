# Builder images

Each subdirectory contains a `Dockerfile` that CI (`build-images.yml`) uses to
build a registry server image. The builder is selected by the `image.builder`
field in the server's manifest; omitting it (the default) selects `go-static`.

## Selecting a builder

In a manifest's `image` block:

```yaml
image:
  ref: ghcr.io/gigmcp/myserver-mcp
  digest: sha256:...
  entrypoint: /app/server
  builder: node          # omit or "go-static" for static Go binaries
```

`registryctl build-args` reads the manifest and emits `BUILDER=<value>` into
`$GITHUB_ENV`; the workflow passes `-f "images/${BUILDER}/Dockerfile"` to
`docker buildx build`.

## The /app/server convention

Every builder places the executable at `/app/server` inside the image. The
manifest's `image.entrypoint` should always be `/app/server`. This is the path
the gateway uses when it mounts or runs the server.

| Builder     | Runtime image base              | /app/server artifact                          |
|-------------|----------------------------------|-----------------------------------------------|
| `go-static` | `scratch` (no OS)                | Static ELF binary (CGO_ENABLED=0)             |
| `toolpack`  | `scratch` (no OS)                | Static toolpack engine + baked-in `/app/manifest.yaml` and `/app/toolspec.yaml` |
| `node`      | `node:22-bookworm-slim`          | esbuild-bundled JS with `#!/usr/bin/env node` shebang |
| `python`    | `python:3.13-slim-bookworm`      | shiv zipapp with `#!/usr/local/bin/python3` shebang   |

## Source conventions per builder

### go-static

- `source.repo` — the Go module root or any git repo with Go source.
- `source.package` — the directory containing the `main` package
  (e.g. a subdirectory of the author's repo); defaults to `.` (repo root).
- The builder runs `CGO_ENABLED=0 go build -trimpath -o /out/server .` in that
  directory.

### toolpack

- `source.repo` / `source.tag` — the generic engine
  (`github.com/gigmcp/toolpack`); built exactly like `go-static`.
- Additionally bakes `manifests/<name>/<version>.yaml` and
  `manifests/<name>/<version>.toolspec.yaml` from this repo (the build context) into
  the image as `/app/manifest.yaml` and `/app/toolspec.yaml`. The engine reads
  the manifest for the credential-inject contract and the toolspec for the
  tool→HTTP mappings; `registryctl lint-toolspecs` enforces their coherence
  in CI.

#### Runtime bundle contract

Unlike single-file (`go-static`) servers, a `toolpack` image is a small
**bundle** the gateway must extract and mount as a unit — the engine
`log.Fatal`s at startup if either sidecar is missing. The contract the
gateway relies on (discriminated by `image.builder == "toolpack"` in the
signed index):

- **Location** — the bundle is exactly `dirname(image.entrypoint)` = `/app`,
  a flat directory. No subdirectories, no files outside it.
- **Contents** — exactly three regular files: the entrypoint binary
  (`/app/server`, mode 0755) plus two inert sidecars `/app/manifest.yaml` and
  `/app/toolspec.yaml` (mode 0644). No symlinks, no nested dirs, no second
  executable. The engine's `GIG_TOOLPACK_MANIFEST` / `GIG_TOOLPACK_SPEC` env
  overrides exist but are stripped by the sandbox's `--clearenv`; the fixed
  `/app` paths are authoritative.
- **Entrypoint** — always `/app/server`; the gateway derives the bundle dir
  from `dirname(image.entrypoint)`.
- **Bounds** — the gateway caps the bundle (≤ 8 files, ≤ 64 MiB total, each
  non-entrypoint file ≤ 1 MiB) and mounts it read-only.
- **Integrity** — `/app/manifest.yaml` is the registry manifest copied
  verbatim, so it is byte-identical to the manifest compiled into the signed
  index; the gateway cross-checks the two (parsed equality) as
  defense-in-depth. `toolspec.yaml` is not in the index — its integrity comes
  from the signed, digest-pinned image.

If this set ever needs to grow (new sidecar, non-flat layout), it is a
breaking change to the bundle contract and must be coordinated with the
gateway before shipping.

### node

- `source.repo` / `source.tag` — a git repo whose checkout contains a
  `package.json` at `source.package` (defaults to repo root).
- The `package.json` must declare either a `bin` entry (string or object with
  one entry) or a `main` entry pointing to the server's JS/TS entry point.
  esbuild bundles everything into a single self-contained CJS file.

### python

- `source.repo` / `source.tag` — a git repo containing a `pyproject.toml` at
  `source.package` (defaults to repo root).
- The `pyproject.toml` must declare **exactly one** `[project.scripts]` entry;
  shiv uses that entry's name as the console-script entrypoint and bundles the
  project plus all dependencies into a single executable zipapp.

## Runnability caveat for node and python

**node and python images are NOT yet installable.** The gigmcp gateway
sandbox mounts a single static binary from a `scratch`-based image. Support
for runtime-rootfs images (full OS layer, node/python interpreter on PATH)
requires the gateway's rootfs sandbox extension, which is designed but not yet
shipped (see `gigmcp docs/superpowers/specs/2026-06-06-rootfs-spike-findings.md`).

These builders exist so manifests can declare `builder: node` or
`builder: python` and CI can build the images today; the images will become
installable once the gateway extension ships.

## Custom per-server Dockerfiles

For servers that don't fit any generic builder (unusual CGO flags, pre-built
assets, multi-step toolchains, etc.) you may add `images/<name>/Dockerfile`.
This overrides the generic builder — point `image.builder` at the directory
name that matches your custom Dockerfile. Custom Dockerfiles are the exception;
prefer the generic builders when possible.
