# Builder for TypeScript/Node MCP servers. Produces a runtime-rootfs image:
# node runtime + a single esbuild-bundled entry at /app/server (node shebang).
# NOTE: rootfs images are NOT yet installable — the gigmcp gateway's v1
# sandbox mounts a single static binary; the rootfs sandbox extension is
# designed (see gigmcp docs/superpowers/specs/2026-06-06-rootfs-spike-findings.md)
# and these images become installable when it ships.
FROM node:22-bookworm-slim AS build
ARG SOURCE_REPO
ARG SOURCE_TAG
ARG PACKAGE=.
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates && rm -rf /var/lib/apt/lists/*
RUN git clone --depth 1 --branch "${SOURCE_TAG}" "https://${SOURCE_REPO}" /repo
WORKDIR /repo/${PACKAGE}
RUN npm ci || npm install
# Bundle to ONE self-contained file so the manifest entrypoint stays /app/server.
# Convention: the package's "main"/"bin" entry is the server; esbuild inlines deps.
RUN npx esbuild "$(node -e "const p=require('./package.json');console.log(p.bin ? (typeof p.bin==='string'?p.bin:Object.values(p.bin)[0]) : p.main || 'index.js')")" \
      --bundle --platform=node --format=cjs --outfile=/out/server.js \
 && printf '#!/usr/bin/env node\n' | cat - /out/server.js > /out/server \
 && chmod 0755 /out/server

FROM node:22-bookworm-slim
COPY --from=build /out/server /app/server
ENTRYPOINT ["/app/server"]
