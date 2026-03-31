#!/usr/bin/env bash
set -euo pipefail

# Clean build artifacts
# Usage: clean.sh <build_dir> <xcodebuild> <project> <scheme> <derived_data>

BUILD_DIR="${1}"
XCODEBUILD="${2:-xcodebuild}"
PROJECT="${3:-ios.xcodeproj}"
SCHEME="${4:-ios}"
DERIVED_DATA="${5}"

rm -rf "${BUILD_DIR}"
${XCODEBUILD} clean \
  -project "${PROJECT}" \
  -scheme "${SCHEME}" \
  -derivedDataPath "${DERIVED_DATA}" \
  -quiet

