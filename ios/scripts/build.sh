#!/usr/bin/env bash
set -euo pipefail

# Build the iOS project
# Usage: build.sh <xcodebuild> <project> <scheme> <destination> <derived_data>

XCODEBUILD="${1:-xcodebuild}"
PROJECT="${2:-ios.xcodeproj}"
SCHEME="${3:-ios}"
DESTINATION="${4}"
DERIVED_DATA="${5}"

${XCODEBUILD} build \
  -project "${PROJECT}" \
  -scheme "${SCHEME}" \
  -destination "${DESTINATION}" \
  -derivedDataPath "${DERIVED_DATA}" \
  -quiet

