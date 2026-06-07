# Builder for Python MCP servers. Produces a runtime-rootfs image: python
# runtime + a single shiv zipapp at /app/server (python shebang).
# NOTE: rootfs images are NOT yet installable — see images/node/Dockerfile note.
FROM python:3.13-slim-bookworm AS build
ARG SOURCE_REPO
ARG SOURCE_TAG
ARG PACKAGE=.
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates && rm -rf /var/lib/apt/lists/*
RUN git clone --depth 1 --branch "${SOURCE_TAG}" "https://${SOURCE_REPO}" /repo
WORKDIR /repo/${PACKAGE}
RUN pip install --no-cache-dir shiv
# Convention: the project is pip-installable (pyproject.toml) and declares ONE
# console script; shiv bundles project+deps into a single executable zipapp.
RUN ENTRY="$(python3 -c "import tomllib, pathlib; data=tomllib.loads(pathlib.Path('pyproject.toml').read_text()); scripts=data.get('project',{}).get('scripts',{}); assert len(scripts)==1, f'need exactly one [project.scripts] entry, got {list(scripts)}'; print(next(iter(scripts)))")" \
 && shiv -c "$ENTRY" -o /out/server --python "/usr/local/bin/python3" .

FROM python:3.13-slim-bookworm
COPY --from=build /out/server /app/server
ENTRYPOINT ["/app/server"]
