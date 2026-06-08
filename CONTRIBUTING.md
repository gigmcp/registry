# Contributing a server

**Aggregator policy:** your server's source code lives in YOUR repo. This
registry holds only your manifest, an optional toolspec (declarative data),
and, rarely, a custom Dockerfile. Do not add source code, `examples/`
directories, or any `.go` files outside `schema/` and `cmd/` — CI will
reject them.

## Two ways to contribute

**Toolspec-driven (no code at all):** if your service's tools are plain REST
calls, you don't need a server repo. Add two files:
`manifests/<name>/<version>.yaml` (with `source.repo:
github.com/gigmcp/toolpack`, `builder: toolpack`) and
`manifests/<name>/<version>.toolspec.yaml` mapping each manifest tool to an
HTTP endpoint (see any existing `*.toolspec.yaml`; `ably` is the canonical
example). `registryctl lint-toolspecs` enforces: tool set matches the
manifest 1:1, every base URL host is within the manifest's egress allowlist,
entrusted tier carries an `auth` block (sealed must not). The generic
[toolpack engine](https://github.com/gigmcp/toolpack) serves the spec; no
server code is written or reviewed.

**Own server (the steps below):** for anything beyond declarative HTTP —
local computation, multi-step flows, non-REST protocols.

1. Tag a release of your MCP server's source repo.
2. Open a PR adding **one file**: `manifests/<name>/<version>.yaml`:
   - `source.repo` / `source.tag` — the tagged source to build from (your repo)
   - `source.package` — subdirectory containing the server's `main` package;
     omit or leave blank to build from the repo root (`.`)
   - `image.builder` — selects the build recipe (`images/<builder>/Dockerfile`);
     omit or set to `go-static` for static Go binaries (the default).
     `node` and `python` builders are **prepared but not yet installable** —
     they require the gigmcp gateway's rootfs sandbox extension (designed,
     pending shipment). See `images/README.md` for details.
   - `image.entrypoint` must be `/app/server` (every builder places the
     server there)
   - egress allowlist: exact hostnames or `*.suffix` (≥2 labels) only
   - credentials: `inject.header`+`format` (sealed) or `inject.env` (entrusted)
   - tools: mark the curated subset `default: true` — ONLY default tools are
     exposed to clients; a manifest with no default tools exposes nothing
3. CI lints the manifest. A maintainer dispatches `build-images`; update your
   PR to pin the printed digest.
4. Owner review + merge ⇒ the signed index is republished automatically.

The generic builder (`images/go-static/Dockerfile`) handles standard static Go
servers. If your server requires an unusual build (custom CGO flags, a
non-Go toolchain, pre-built assets, etc.) you may add
`images/<name>/Dockerfile` — but this is the exception, not the rule.

Manifest changes on version bump force re-consent in every gateway that has
the server installed.
