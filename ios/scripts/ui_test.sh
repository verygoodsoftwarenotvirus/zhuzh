#!/usr/bin/env bash
set -euo pipefail

# Run UI tests
# Usage: ui_test.sh <xcodebuild> <project> <scheme> <destination> <ui_test_target> <derived_data>

XCODEBUILD="${1:-xcodebuild}"
PROJECT="${2:-ios.xcodeproj}"
SCHEME="${3:-ios}"
DESTINATION="${4}"
UI_TEST_TARGET="${5:-iosUITests}"
DERIVED_DATA="${6}"

${XCODEBUILD} test \
  -project "${PROJECT}" \
  -scheme "${SCHEME}" \
  -destination "${DESTINATION}" \
  -only-testing:"${UI_TEST_TARGET}" \
  -derivedDataPath "${DERIVED_DATA}" \
  -quiet

