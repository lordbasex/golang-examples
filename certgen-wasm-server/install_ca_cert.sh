#!/usr/bin/env bash
set -euo pipefail

# -----------------------------------------------------------------------------
# Cross-platform CA installer (macOS, Linux, Windows via PowerShell).
# Usage:
#   ./install_ca_cert.sh
# The script always uses ca.crt in the same directory as this script.
# -----------------------------------------------------------------------------

log() {
  printf '[INFO] %s\n' "$1"
}

err() {
  printf '[ERROR] %s\n' "$1" >&2
}

require_command() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    err "Required command not found: $cmd"
    exit 1
  fi
}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CERT_INPUT="$SCRIPT_DIR/ca.crt"

CERT_PATH="$(cd "$(dirname "$CERT_INPUT")" && pwd)/$(basename "$CERT_INPUT")"
CERT_FILE_NAME="$(basename "$CERT_PATH")"

if [[ ! -f "$CERT_PATH" ]]; then
  err "Certificate file does not exist: $CERT_PATH"
  exit 1
fi

require_command openssl

# Convert fingerprint to the same format used by Windows Thumbprint.
CERT_SHA256_HEX="$(
  openssl x509 -in "$CERT_PATH" -noout -fingerprint -sha256 \
  | awk -F= '{print $2}' \
  | tr -d ':' \
  | tr '[:lower:]' '[:upper:]'
)"

detect_os() {
  local uname_out
  uname_out="$(uname -s 2>/dev/null || true)"

  case "${uname_out}" in
    Darwin*) echo "macos" ;;
    Linux*) echo "linux" ;;
    CYGWIN*|MINGW*|MSYS*) echo "windows" ;;
    *)
      if [[ "${OS:-}" == "Windows_NT" ]]; then
        echo "windows"
      else
        echo "unknown"
      fi
      ;;
  esac
}

ROOT_CMD=()

init_privilege_mode() {
  local os_name="$1"
  if [[ "$os_name" == "windows" ]]; then
    return 0
  fi

  if [[ "${EUID:-$(id -u)}" -eq 0 ]]; then
    ROOT_CMD=()
    return 0
  fi

  require_command sudo

  # Non-interactive sudo: if password is required, command fails immediately.
  if ! sudo -n true >/dev/null 2>&1; then
    err "This script runs in silent/non-interactive mode."
    err "Grant admin privileges first (root user or sudo without password prompt)."
    exit 1
  fi

  ROOT_CMD=(sudo -n)
}

run_as_root() {
  if [[ ${#ROOT_CMD[@]} -gt 0 ]]; then
    "${ROOT_CMD[@]}" "$@"
  else
    "$@"
  fi
}

install_macos() {
  log "Detected macOS. Installing certificate into System keychain (silent mode)..."
  run_as_root security add-trusted-cert \
    -d \
    -r trustRoot \
    -k /Library/Keychains/System.keychain \
    "$CERT_PATH" >/dev/null 2>&1

  # Deterministic validation: compare SHA-256 hash format from security -Z.
  local fp_hex
  fp_hex="$CERT_SHA256_HEX"

  log "Validating certificate installation on macOS..."
  if run_as_root security find-certificate -a -Z /Library/Keychains/System.keychain | awk '{print toupper($0)}' | awk -v fp="$fp_hex" 'index($0, fp){found=1} END{exit(found?0:1)}'; then
    log "Validation successful on macOS."
  else
    err "Validation failed on macOS. Fingerprint not found in System keychain."
    exit 1
  fi
}

install_linux() {
  log "Detected Linux. Installing certificate into system CA store (silent mode)..."
  local ca_bundle=""
  local installed_path=""

  if command -v update-ca-certificates >/dev/null 2>&1; then
    # Debian/Ubuntu/Alpine-like flow.
    installed_path="/usr/local/share/ca-certificates/${CERT_FILE_NAME}"
    run_as_root cp "$CERT_PATH" "$installed_path"
    run_as_root update-ca-certificates >/dev/null 2>&1
    ca_bundle="/etc/ssl/certs/ca-certificates.crt"
  elif command -v update-ca-trust >/dev/null 2>&1; then
    # RHEL/CentOS/Fedora-like flow.
    installed_path="/etc/pki/ca-trust/source/anchors/${CERT_FILE_NAME}"
    run_as_root cp "$CERT_PATH" "$installed_path"
    run_as_root update-ca-trust extract >/dev/null 2>&1
    ca_bundle="/etc/pki/tls/certs/ca-bundle.crt"
  else
    err "Could not detect Linux CA update tool (update-ca-certificates or update-ca-trust)."
    err "Install manually for your distro."
    exit 1
  fi

  if [[ ! -f "$ca_bundle" ]]; then
    err "CA bundle not found after installation: $ca_bundle"
    exit 1
  fi

  # Validation without chain checks: verify file copied and fingerprint preserved.
  if [[ ! -f "$installed_path" ]]; then
    err "Validation failed on Linux. Installed certificate not found: $installed_path"
    exit 1
  fi

  local installed_sha
  installed_sha="$(
    openssl x509 -in "$installed_path" -noout -fingerprint -sha256 \
    | awk -F= '{print $2}' \
    | tr -d ':' \
    | tr '[:lower:]' '[:upper:]'
  )"

  log "Validating certificate installation on Linux..."
  if [[ "$installed_sha" == "$CERT_SHA256_HEX" ]]; then
    log "Validation successful on Linux."
  else
    err "Validation failed on Linux. Installed certificate fingerprint mismatch."
    exit 1
  fi
}

install_windows() {
  log "Detected Windows. Installing certificate into LocalMachine Root store (silent mode)..."
  local ps_cmd=""
  local cert_for_windows="$CERT_PATH"

  if command -v powershell.exe >/dev/null 2>&1; then
    ps_cmd="powershell.exe"
  elif command -v pwsh >/dev/null 2>&1; then
    ps_cmd="pwsh"
  else
    err "PowerShell not found (powershell.exe or pwsh)."
    exit 1
  fi

  # Convert to Windows path when running from Git Bash/MSYS/Cygwin.
  if command -v cygpath >/dev/null 2>&1; then
    cert_for_windows="$(cygpath -w "$CERT_PATH")"
  fi

  "$ps_cmd" -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command \
    "\$ErrorActionPreference='Stop'; Import-Certificate -FilePath '$cert_for_windows' -CertStoreLocation 'Cert:\LocalMachine\Root' | Out-Null"

  log "Validating certificate installation on Windows..."
  "$ps_cmd" -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command \
    "\$ErrorActionPreference='Stop'; \$thumb='$CERT_SHA256_HEX'; \$cert = Get-ChildItem -Path 'Cert:\LocalMachine\Root' | Where-Object { \$_.Thumbprint -eq \$thumb }; if (-not \$cert) { throw 'Certificate not found in LocalMachine\Root.' }"

  log "Validation successful on Windows."
}

main() {
  local os_name
  os_name="$(detect_os)"
  init_privilege_mode "$os_name"

  case "$os_name" in
    macos) install_macos ;;
    linux) install_linux ;;
    windows) install_windows ;;
    *)
      err "Unsupported operating system."
      err "Detected uname: $(uname -s 2>/dev/null || echo unknown)"
      exit 1
      ;;
  esac

  log "Done. Certificate installed and validated."
}

main
