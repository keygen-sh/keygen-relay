#!/bin/bash

log_info() {
  echo "[info] $1"
}

log_err() {
  echo "[error] $1"
  exit 1
}

main() {
  make build-all # build cross-platform

  # copy installer and version to dist
  cp ./scripts/install.sh ./dist/install.sh
  cp ./VERSION ./dist/version
}

main
