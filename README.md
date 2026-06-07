# registry

Curated, signed registry of community MCP servers for [gigmcp](https://github.com/gigmcp/gigmcp).
Servers are packaged as digest-pinned OCI images with author-declared
entitlement manifests (egress allowlist, credential schema, tool subset,
security tier). The gateway installs from the **signed `index.json`** release
artifact — never from this repo directly.

## Catalog

The registry tracks 224 manifests mirroring Composio's toolkit list. `github`,
`echo`, and `fetch` have reference implementations in their own repos (see
CATALOG.md), while the rest are **planned** — curated, lint-enforced manifests
whose implementations and image digests are pending (placeholder digests are not
installable). See [CATALOG.md](CATALOG.md) for the full table of names, auth
types, tiers, and egress allowlists.

**Aggregator policy:** this repository holds manifests and build recipes only.
Server source code always lives in the author's own repo; it is never committed
here.

## Layout
- `manifests/<name>/<version>.yaml` — one manifest per server version (schema: `schema/`)
- `schema/` — Go module (Apache-2.0): the authoritative parser/validator used
  byte-for-byte by CI here and by the gigmcp gateway
- `denylist/exfil-domains.txt` — egress domains lint CI rejects
- `images/go-static/Dockerfile` — generic static-binary builder (v1: FROM scratch, static ELF); `images/node/Dockerfile` and `images/python/Dockerfile` — prepared runtime-rootfs builders (NOT yet installable — pending gateway rootfs sandbox extension); custom `images/<name>/Dockerfile` for unusual builds; see `images/README.md`
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
(the current echo/fetch ones) are not installable until CI has built their
images and the digests are pinned.

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

`echo` and `fetch` reference the gigmcp repo as source and become buildable
once that repo is public.
