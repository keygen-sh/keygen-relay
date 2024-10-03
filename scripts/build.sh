#!/bin/bash

export CGO_ENABLED=0

log_info() {
  echo "[info] $1"
}

log_err() {
  echo "[error] $1"
  exit 1
}

main() {
  # Clear out any previous builds
  rm build/*

  for platform in $PLATFORMS
  do
    IFS='/' read -r os arch <<< "$platform"

    filename="relay_${os}_${arch}"
    if [ "${os}" = 'windows' ]
    then
      filename="${filename}.exe"
    fi

    log_info "building v${VERSION} for ${platform}"

    env GOOS="${os}" GOARCH="${arch}" \
      go build -o "build/${filename}" -ldflags "-X ${PACKAGE}.Version=${VERSION}" ./cmd/relay

    if [ $? -eq 0 ]
    then
      log_info "successfully built v${VERSION} for ${platform}"
    else
      log_err "failed to build v${VERSION} for ${platform}"
    fi
  done

  # Copy installer and version to build dir
  cp ./scripts/install.sh ./build/install.sh
  cp ./VERSION ./build/version
}

# FIXME(ezekg) Cross-compiling these distros on darwin/amd64 fails
PLATFORMS="$(go tool dist list | grep -vE 'ios|android|js|aix|illumos|riscv64|plan9|solaris|loong')"
PACKAGE='github.com/keygen-sh/keygen-relay/cli'
VERSION="$(cat VERSION)"

main
