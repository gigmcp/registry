# Generic builder for static Go MCP servers (v1 packaging). Build args come
# from the server's manifest (source.repo, source.tag, source.package) via
# `registryctl build-args` — see .github/workflows/build-images.yml.
# Servers needing a custom build may still add images/<name>/Dockerfile.
FROM golang:1.26-bookworm AS build
ARG SOURCE_REPO
ARG SOURCE_TAG
ARG PACKAGE=.
RUN git clone --depth 1 --branch "${SOURCE_TAG}" "https://${SOURCE_REPO}" /repo
WORKDIR /repo/${PACKAGE}
RUN CGO_ENABLED=0 go build -trimpath -o /out/server .

FROM scratch
COPY --from=build /out/server /app/server
ENTRYPOINT ["/app/server"]
