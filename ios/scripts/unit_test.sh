#!/usr/bin/env bash
set -euo pipefail

# Run unit tests
# Usage: unit_test.sh <xcodebuild> <project> <scheme> <destination> <test_target> <derived_data>

XCODEBUILD="${1:-xcodebuild}"
PROJECT="${2:-ios.xcodeproj}"
SCHEME="${3:-ios}"
DESTINATION="${4}"
TEST_TARGET="${5:-iosTests}"
DERIVED_DATA="${6}"

${XCODEBUILD} test \
  -project "${PROJECT}" \
  -scheme "${SCHEME}" \
  -destination "${DESTINATION}" \
  -only-testing:"${TEST_TARGET}" \
  -derivedDataPath "${DERIVED_DATA}" \
  -parallel-testing-enabled NO \
  -quiet

