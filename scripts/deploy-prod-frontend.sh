#!/usr/bin/env bash
# Run the prod frontend Skaffold deploy locally (consumer + admin webapps only).
# Prerequisites: Docker, gcloud auth, kubectl pointed at prod cluster.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

echo "=== Prerequisites check ==="
command -v docker >/dev/null || { echo "Docker required"; exit 1; }
command -v gcloud >/dev/null || { echo "gcloud required"; exit 1; }
command -v kubectl >/dev/null || { echo "kubectl required"; exit 1; }
command -v skaffold >/dev/null || { echo "skaffold required (gcloud components install skaffold)"; exit 1; }

echo "=== Fetching GKE credentials ==="
gcloud container clusters get-credentials prod --region us-central1

echo "=== Running Skaffold deploy (frontend only, profile prod) ==="
cd frontend
VERSION="${VERSION:-local}" skaffold run \
  --filename=skaffold.yaml \
  --build-concurrency 1 \
  --profile prod \
  --label "deploy-source=local" \
  "$@"
