#!/bin/sh

log_info() {
  echo "[info] ${1}"
}

log_err() {
  echo "[error] ${1} (please install manually via https://keygen.sh/docs/relay/)"
  exit 1
}

get_os() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')

  case "${os}"
  in
  msys*|mingw*)
    os='windows'
    ;;
  cygwin*)
    os='linux'
    ;;
  esac

  if [ -z "${os}" ]
  then
    log_err 'unable to detect operating system'
  fi

  echo "${os}"
}

get_arch() {
  arch=$(uname -m)

  case "${arch}"
  in
  x86_64|amd64p32)
    arch='amd64'
    ;;
  x86)
    arch='386'
    ;;
  aarch64)
    arch='arm64'
    ;;
  armv8)
    arch='arm64'
    ;;
  armv*)
    arch='arm'
    ;;
  i686|i386)
    arch='386'
    ;;
  esac

  if [ -z "${arch}" ]
  then
    log_err 'unable to detect architecture'
  fi

  echo "${arch}"
}

get_tmp_path() {
  path=$(mktemp -d /tmp/relay.XXXXXXXXXX)
  if [ -z "${path}" ]
  then
    log_err 'unable to make tmp install path'
  fi

  echo "${path}/${BIN_NAME}"
}

get_bin_name() {
  name='relay'
  if [ "${OS}" = 'windows' ]
  then
    name="${name}.exe"
  fi

  echo "${name}"
}

get_bin_path() {
  path="/usr/local/bin/${BIN_NAME}"
  if [ "${OS}" = 'windows' ]
  then
    path="./${BIN_NAME}"
  fi

  echo "${path}"
}

get_bin_version() {
  version=$(curl -sSL 'https://get.keygen.sh/keygen/relay/latest/version')
  if [ -z "${version}" ]
  then
    log_err 'unable to get latest version'
  fi

  echo "${version}"
}

get_bin_url() {
  filename="relay_${OS}_${ARCH}"
  if [ "${OS}" = 'windows' ]
  then
    filename="${filename}.exe"
  fi

  echo "https://get.keygen.sh/keygen/relay/${BIN_VERSION}/${filename}"
}

assert_os_support() {
  case "${OS}"
  in
    linux) return ;;
  esac

  log_err "unsupported operating system: ${OS}"
}

assert_arch_support() {
  case "${ARCH}"
  in
    386)   return ;;
    amd64) return ;;
    arm64) return ;;
    arm)   return ;;
  esac

  log_err "unsupported architecture: ${ARCH}"
}

assert_platform_support() {
  assert_os_support
  assert_arch_support
}

main() {
  assert_platform_support

  status=$(curl -sSL "${BIN_URL}" --write-out "%{http_code}" -o "${BIN_TMP}")

  if [ "${status}" -eq 200 ]
  then
    log_info "successfully downloaded v${BIN_VERSION} for ${PLATFORM}: ${BIN_TMP}"
  else
    log_err "failed to download v${BIN_VERSION} for ${PLATFORM}: ${status}"
  fi

  mv "${BIN_TMP}" "${BIN_PATH}" && \
    chmod +x "${BIN_PATH}"

  if [ $? -eq 0 ]
  then
    log_info "successfully installed v${BIN_VERSION} for ${PLATFORM}: ${BIN_PATH}"
  else
    log_err "failed to installed v${BIN_VERSION} for ${PLATFORM}"
  fi

  ${BIN_PATH} --help
}

OS=$(get_os)
ARCH=$(get_arch)
PLATFORM="${OS}/${ARCH}"
BIN_VERSION=$(get_bin_version)
BIN_NAME=$(get_bin_name)
BIN_PATH=$(get_bin_path)
BIN_TMP=$(get_tmp_path)
BIN_URL=$(get_bin_url)

main
