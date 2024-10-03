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

  # draft new version
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

  # upload artifact for each dist platform
  for filename in $(ls dist/relay_*)
  do
    IFS='_' read -r _ platform arch <<< "$filename"

    log_info "uploading v${VERSION} for ${platform}/${arch}: ${filename}"

    keygen upload "build/${filename}" \
      --release "${VERSION}" \
      --platform "${platform}" \
      --arch "${arch}"

    if [ $? -eq 0 ]
    then
      log_info "successfully uploading v${VERSION} for ${platform}/${arch}"
    else
      log_err "failed to upload v${VERSION} for ${platform}/${arch}"
    fi
  done

  # upload installer
  keygen upload 'build/install.sh' --release "${VERSION}"

  # upload version
  keygen upload 'build/version' \
    --release "${VERSION}" \
    --filetype 'txt'

  # publish version
  keygen publish --release "${VERSION}"

  if [ $? -eq 0 ]
  then
    log_info "successfully published v${VERSION}"
  else
    log_err "failed to publish v${VERSION}"
  fi

  # we only want to do the rest for stable releases
  if [ "${CHANNEL}" = 'stable' ]
  then
    # untag previous latest if it exists (we'll continue even on failure)
    keygen untag --release 'latest'

    if [ $? -eq 0 ]
    then
      log_info "successfully untagged v${VERSION}"
    else
      log_warn "failed to untag v${VERSION}"
    fi

    # tag as latest
    keygen tag 'latest' --release "${VERSION}"

    if [ $? -eq 0 ]
    then
      log_info "successfully tagged v${VERSION}"
    else
      log_err "failed to tag v${VERSION}"
    fi
  fi
}

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
