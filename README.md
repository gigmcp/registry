# registry

Curated, signed registry of community MCP servers for [gigmcp](https://github.com/gigmcp/gigmcp).
Servers are packaged as digest-pinned OCI images with author-declared
entitlement manifests (egress allowlist, credential schema, tool subset,
security tier). The gateway installs from the **signed `index.json`** release
artifact — never from this repo directly.

## Catalog

The registry tracks 221 manifests mirroring Composio's toolkit list. Most are
**toolspec-driven**: `manifests/<name>/<version>.toolspec.yaml` maps each manifest
tool to a real REST endpoint, served by the generic
[toolpack](https://github.com/gigmcp/toolpack) engine and built by the
`toolpack` builder. A few adopt an established upstream Go MCP server
instead, and a handful with no usable public API remain planned. Image
digests are placeholders (`sha256:0000…`) until `build-images` CI pins the
real ones — placeholder digests are not installable. See
[CATALOG.md](CATALOG.md) for the full table of names, auth types, tiers,
egress allowlists, and per-entry status.

**Aggregator policy:** this repository holds manifests and build recipes only.
Server source code always lives in the author's own repo; it is never committed
here.

## Layout
- `manifests/<name>/<version>.yaml` — one manifest per server version (schema: `schema/`)
- `manifests/<name>/<version>.toolspec.yaml` — declarative tool→HTTP mapping
  consumed by the generic [toolpack](https://github.com/gigmcp/toolpack) engine
  (manifests with `builder: toolpack`); sits beside its manifest version and is
  lint-enforced against it (`registryctl lint-toolspecs`)
- `schema/` — Go module (Apache-2.0): the authoritative parser/validator used
  byte-for-byte by CI here and by the gigmcp gateway
- `denylist/exfil-domains.txt` — egress domains lint CI rejects
- `images/go-static/Dockerfile` — generic static-binary builder (FROM scratch, static ELF); `images/toolpack/Dockerfile` — toolpack-engine builder (static engine + baked-in manifest/toolspec); `images/node/Dockerfile` and `images/python/Dockerfile` — prepared runtime-rootfs builders (NOT yet installable — pending gateway rootfs sandbox extension); custom `images/<name>/Dockerfile` for unusual builds; see `images/README.md`
- `cmd/registryctl` — lint | build-index | sign | verify | keygen

## Trust chain
1. Manifests are PR-gated; lint CI blocks invalid schemas, broad wildcards,
   and denylisted exfil domains. The path `manifests/<name>/<version>.yaml`
   must match the manifest contents.
2. Images are built by CI from the author's tagged source; the manifest pins
   the resulting **linux/amd64 image-manifest digest** — what was approved is
   what runs.
3. On merge to main, CI compiles all manifests into `index.json`, signs it
   (ed25519), and publishes both as the rolling `latest` release. The gateway
   verifies the signature before trusting any entry.

## Bootstrapping a manifest's digest
`build-images.yml` (manual dispatch) builds and pushes the image and prints
the digest; the PR then pins that digest. Manifests with placeholder digests
are not installable until CI has built their images and the digests are
pinned.

## Signing key
Generate once with `go run ./cmd/registryctl keygen`. Private key → repo
secret `GIG_INDEX_SIGNING_KEY`. Public key → gateway `GIG_REGISTRY_PUBKEY`.

## Bootstrap

1. Push this repo; `lint` CI must be green.
2. Generate keys (`go run ./cmd/registryctl keygen`); set the private key as
   repo secret `GIG_INDEX_SIGNING_KEY`.
3. `publish-index` signs and releases `index.json` + `index.json.sig` on every
   push to main.
4. Making a catalog entry installable: publish the server's source repo and
   tag it, dispatch `build-images` with its name (version optional, defaults
   to latest), pin the printed digest in `manifests/<name>/<version>.yaml`,
   and merge.
5. Point a gateway at the index: `GIG_REGISTRY_INDEX_URL=<release asset URL>`,
   `GIG_REGISTRY_PUBKEY=<public key>`, `GIG_INSTALL=<name>` — sealed-tier
   servers only ever see a placeholder token; the egress proxy injects the
   real credential for the manifest's allowlisted hosts only.
