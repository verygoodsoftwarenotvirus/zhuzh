#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=branding.sh
source "$(dirname "${BASH_SOURCE[0]}")/branding.sh"

# Proxy Cloud SQL database
# Usage: proxy_db.sh [gcp_project_name] [instance_name]

GCP_PROJECT_NAME="${1:-$GCP_PROJECT_ID}"
INSTANCE_NAME="${2:-prod}"

cloud_sql_proxy "${GCP_PROJECT_NAME}:us-central1:${INSTANCE_NAME}" --port 5434 --gcloud-auth
