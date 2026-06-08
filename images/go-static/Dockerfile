# Generic builder for static Go MCP servers. Build args come
# from the server's manifest (source.repo, source.tag, source.package) via
# `registryctl build-args` — see .github/workflows/build-images.yml.
# Servers needing a custom build may still add images/<name>/Dockerfile.
#
# Multi-arch: the build stage is pinned to $BUILDPLATFORM (runs natively on the
# runner) and cross-compiles to $TARGETOS/$TARGETARCH, so `docker buildx
# --platform linux/amd64,linux/arm64` produces one static binary per target
# with no QEMU — the resulting OCI image index resolves per host architecture.
FROM --platform=$BUILDPLATFORM golang:1.26-bookworm AS build
ARG SOURCE_REPO
ARG SOURCE_TAG
ARG PACKAGE=.
ARG TARGETOS
ARG TARGETARCH
RUN git clone --depth 1 --branch "${SOURCE_TAG}" "https://${SOURCE_REPO}" /repo
WORKDIR /repo/${PACKAGE}
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -o /out/server .

FROM scratch
COPY --from=build /out/server /app/server
ENTRYPOINT ["/app/server"]
