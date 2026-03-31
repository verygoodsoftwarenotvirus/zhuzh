#!/usr/bin/env bash
set -euo pipefail

# Lint markdown files
# Usage: lint_markdown.sh <container_runner>

CONTAINER_RUNNER="${1:-docker}"

"${CONTAINER_RUNNER}" run --rm --volume $PWD:$PWD --workdir=$PWD --user $(id -u):$(id -g) ghcr.io/igorshubovych/markdownlint-cli:latest --ignore "**/vendor/**" --ignore "**/node_modules/**" --ignore "**/ios/build/**" "**/*.md" --fix --disable=MD013
