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
| `node`      | `node:22-bookworm-slim`          | esbuild-bundled JS with `#!/usr/bin/env node` shebang |
| `python`    | `python:3.13-slim-bookworm`      | shiv zipapp with `#!/usr/local/bin/python3` shebang   |

## Source conventions per builder

### go-static

- `source.repo` — the Go module root or any git repo with Go source.
- `source.package` — the directory containing the `main` package
  (e.g. a subdirectory of the author's repo); defaults to `.` (repo root).
- The builder runs `CGO_ENABLED=0 go build -trimpath -o /out/server .` in that
  directory.

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

**node and python images are NOT yet installable.** The gigmcp gateway v1
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
