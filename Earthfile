VERSION 0.8

ARG --global TARGET_DOCKER_REGISTRY=ghcr.io/librepod

# See Bun docker documentation: https://bun.sh/guides/ecosystem/docker

# use the official Bun image
# see all versions at https://hub.docker.com/r/oven/bun/tags
FROM oven/bun:1
WORKDIR /usr/src/app

install:
  # install dependencies into temp directory
  # this will cache them and speed up future builds
  RUN mkdir -p /temp/dev
  COPY package.json bun.lockb /temp/dev/
  RUN cd /temp/dev && bun install --frozen-lockfile

  # install with --production (exclude devDependencies)
  RUN mkdir -p /temp/prod
  COPY package.json bun.lockb /temp/prod/
  RUN cd /temp/prod && bun install --frozen-lockfile --production

validate-pr:
  FROM +install
  COPY . .
  RUN cp -r /temp/dev/node_modules node_modules
  RUN bun run lint
  RUN bun run format:check
  RUN bun test

build:
  FROM +install
  ENV NODE_ENV=production
  # copy node_modules from temp directory
  # then copy all (non-ignored) project files into the image
  RUN cp -r /temp/dev/node_modules node_modules
  COPY . .
  RUN bun test
  RUN bun run build
  SAVE ARTIFACT build /build AS LOCAL ./build

image:
  FROM oven/bun:1
  LABEL org.opencontainers.image.source="https://github.com/librepod/frp-port-keeper"
  ARG RELEASE_VERSION=latest
  ENV NODE_ENV=production
  ENV ALLOW_PORTS=8000-29999
  COPY +build/build /app/build
  RUN ls -al /app/build
  ENTRYPOINT ["/app/build/frp-port-keeper"]
  SAVE IMAGE --push ${TARGET_DOCKER_REGISTRY}/frp-port-keeper:$RELEASE_VERSION

multi-image:
  BUILD --platform=linux/amd64 --platform=linux/arm64 +build

multi-build:
  FROM +install
  ARG RELEASE_VERSION=latest
  ENV NODE_ENV=production
  # copy node_modules from temp directory
  # then copy all (non-ignored) project files into the image
  RUN cp -r /temp/dev/node_modules node_modules
  COPY . .
  RUN bun build --compile --target=bun-darwin-x64 src/server.ts --outfile build/${RELEASE_VERSION}/darwin/amd64/frp-port-keeper
  RUN bun build --compile --target=bun-darwin-arm64 src/server.ts --outfile build/${RELEASE_VERSION}/darwin/arm64/frp-port-keeper
  RUN bun build --compile --target=bun-windows-x64 src/server.ts --outfile build/${RELEASE_VERSION}/windows/amd64/frp-port-keeper
  RUN bun build --compile --target=bun-linux-x64 src/server.ts --outfile build/${RELEASE_VERSION}/linux/amd64/frp-port-keeper
  RUN bun build --compile --target=bun-linux-arm64 src/server.ts --outfile build/${RELEASE_VERSION}/linux/arm64/frp-port-keeper
  SAVE ARTIFACT build /build AS LOCAL ./build

release:
  FROM  +multi-build

  ARG --required GITHUB_TOKEN
  ARG --required RELEASE_VERSION
  ARG OUT_BASE="./build/${RELEASE_VERSION}"

  # Install gh-cli
  RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg \
      && chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg \
      && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
      && apt-get update && apt-get install gh jq -y \
      && gh --version

  # Generate release notes
  RUN gh api -X POST 'repos/librepod/frp-port-keeper/releases/generate-notes' \
        -F commitish=${RELEASE_VERSION} \
        -F tag_name=${RELEASE_VERSION} \
      > tmp-release-notes.json \
      && cat ./tmp-release-notes.json | jq

  # Gzip the bins
  RUN tar -czvf "${OUT_BASE}/darwin/amd64/frp-port-keeper_darwin_amd64.tar.gz" -C "${OUT_BASE}/darwin/amd64" frp-port-keeper \
      && tar -czvf "${OUT_BASE}/darwin/arm64/frp-port-keeper_darwin_arm64.tar.gz" -C "${OUT_BASE}/darwin/arm64" frp-port-keeper \
      && tar -czvf "${OUT_BASE}/windows/amd64/frp-port-keeper_windows_amd64.tar.gz" -C "${OUT_BASE}/windows/amd64" frp-port-keeper.exe \
      && tar -czvf "${OUT_BASE}/linux/amd64/frp-port-keeper_linux_amd64.tar.gz" -C "${OUT_BASE}/linux/amd64" frp-port-keeper \
      && tar -czvf "${OUT_BASE}/linux/arm64/frp-port-keeper_linux_arm64.tar.gz" -C "${OUT_BASE}/linux/arm64" frp-port-keeper

  # Create release
  RUN jq -r .body tmp-release-notes.json > tmp-release-notes.md \
      && gh release create ${RELEASE_VERSION} \
        --title "$(jq -r .name tmp-release-notes.json)" \
        --notes-file tmp-release-notes.md \
        --verify-tag \
        --draft \
        "${OUT_BASE}/darwin/amd64/frp-port-keeper_darwin_amd64.tar.gz#frp-port-keeper_osx_amd64" \
        "${OUT_BASE}/darwin/arm64/frp-port-keeper_darwin_arm64.tar.gz#frp-port-keeper_osx_arm64" \
        "${OUT_BASE}/windows/amd64/frp-port-keeper_windows_amd64.tar.gz#frp-port-keeper_windows_amd64" \
        "${OUT_BASE}/linux/amd64/frp-port-keeper_linux_amd64.tar.gz#frp-port-keeper_linux_amd64" \
        "${OUT_BASE}/linux/arm64/frp-port-keeper_linux_arm64.tar.gz#frp-port-keeper_linux_arm64"
  
  SAVE ARTIFACT build /build AS LOCAL ./build
