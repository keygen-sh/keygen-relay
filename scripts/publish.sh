#!/bin/bash

log_info() {
  echo "[info] $1"
}

log_warn() {
  echo "[warn] $1"
}

log_err() {
  echo "[error] $1"
  exit 1
}

main() {
  log_info "drafting v${VERSION}"

  # Draft new version
  keygen new \
    --name "Keygen Relay v${VERSION}" \
    --channel "${CHANNEL}" \
    --version "${VERSION}"

  if [ $? -eq 0 ]
  then
    log_info "successfully drafted v${VERSION}"
  else
    log_err "failed to draft v${VERSION}"
  fi

  # Upload artifact for each platform
  for platform in $PLATFORMS
  do
    IFS='/' read -r os arch <<< "$platform"

    filename="relay_${os}_${arch}"
    if [ "${os}" = 'windows' ]
    then
      filename="${filename}.exe"
    fi

    log_info "uploading v${VERSION} for ${platform}: ${filename}"

    keygen upload "build/${filename}" \
      --release "${VERSION}" \
      --platform "${os}" \
      --arch "${arch}"

    if [ $? -eq 0 ]
    then
      log_info "successfully uploading v${VERSION} for ${platform}"
    else
      log_err "failed to upload v${VERSION} for ${platform}"
    fi
  done

  # Upload installer
  keygen upload 'build/install.sh' --release "${VERSION}"

  # Upload version
  keygen upload 'build/version' \
    --release "${VERSION}" \
    --filetype 'txt'

  # Publish version
  keygen publish --release "${VERSION}"

  if [ $? -eq 0 ]
  then
    log_info "successfully published v${VERSION}"
  else
    log_err "failed to publish v${VERSION}"
  fi

  # We only want to do the rest for stable releases
  if [ "${CHANNEL}" = 'stable' ]
  then
    # Untag previous latest if it exists (we'll continue even on failure)
    keygen untag --release 'latest'

    if [ $? -eq 0 ]
    then
      log_info "successfully untagged v${VERSION}"
    else
      log_warn "failed to untag v${VERSION}"
    fi

    # Tag as latest
    keygen tag 'latest' --release "${VERSION}"

    if [ $? -eq 0 ]
    then
      log_info "successfully tagged v${VERSION}"
    else
      log_err "failed to tag v${VERSION}"
    fi
  fi
}

PLATFORMS="$(go tool dist list | grep -vE 'ios|android|js|aix|illumos|riscv64|plan9|solaris|loong')"
PACKAGE='github.com/keygen-sh/keygen-relay/cli'
VERSION="$(cat VERSION)"
CHANNEL='stable'

case "${VERSION}"
in
  *-rc.*)
    CHANNEL='rc'
    ;;
  *-beta.*)
    CHANNEL='beta'
    ;;
  *-alpha.*)
    CHANNEL='alpha'
    ;;
  *-dev.*)
    CHANNEL='dev'
    ;;
esac

main
