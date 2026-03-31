#!/usr/bin/env bash
set -euo pipefail

# Format Swift files using swift-format
# Usage: format.sh <ios_dir> [--check]

IOS_DIR="${1:-ios}"
CHECK_MODE="${2:-}"

if [ "${CHECK_MODE}" = "--check" ]; then
  find "${IOS_DIR}" -name "*.swift" -type f -exec swift-format lint --strict {} +
else
  find "${IOS_DIR}" -name "*.swift" -type f -exec swift-format --in-place {} +
fi

