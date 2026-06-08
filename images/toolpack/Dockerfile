# Builder for toolpack-driven servers: builds the generic toolpack engine
# (github.com/gigmcp/toolpack) from its tagged source, then bakes in the
# server's manifest and toolspec from THIS repo (the build context), so the
# image is fully self-describing. Build args come from `registryctl
# build-args` — see .github/workflows/build-images.yml.
FROM golang:1.26-bookworm AS build
ARG SOURCE_REPO
ARG SOURCE_TAG
ARG PACKAGE=.
RUN git clone --depth 1 --branch "${SOURCE_TAG}" "https://${SOURCE_REPO}" /repo
WORKDIR /repo/${PACKAGE}
RUN CGO_ENABLED=0 go build -trimpath -o /out/server .

FROM scratch
ARG NAME
ARG VERSION
COPY --from=build /out/server /app/server
# The engine reads both at startup: the manifest for the credential inject
# contract, the toolspec for the tool→HTTP mappings.
COPY manifests/${NAME}/${VERSION}.yaml /app/manifest.yaml
COPY manifests/${NAME}/${VERSION}.toolspec.yaml /app/toolspec.yaml
ENTRYPOINT ["/app/server"]
