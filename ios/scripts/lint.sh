#!/usr/bin/env bash
set -euo pipefail

# Lint Swift code using swiftlint
# Usage: lint.sh [--fix]

if [ "${1:-}" = "--fix" ]; then
  swiftlint lint --fix
else
  swiftlint lint --strict
fi

