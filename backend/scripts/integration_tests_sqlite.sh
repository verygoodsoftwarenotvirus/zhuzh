#!/usr/bin/env bash
set -euo pipefail

# Run integration tests for SQLite
# Usage: integration_tests_sqlite.sh <package_prefix>

PACKAGE_PREFIX="${1:-github.com/verygoodsoftwarenotvirus/zhuzh/backend}"

ZHUZH_INTEGRATION_TEST_DB=sqlite go test -v -count=1 "${PACKAGE_PREFIX}/testing/integration/apiserver"
