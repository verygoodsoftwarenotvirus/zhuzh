#!/usr/bin/env bash
# Run terraform plan/apply for prod backend.
# Prerequisites: terraform, gcloud, Terraform Cloud login (terraform login), workspace variables set.
# Variables are typically configured in Terraform Cloud: Workspace → Variables.
#
# Uses token-based kubeconfig (like CI's get-gke-credentials) so gke-gcloud-auth-plugin is not required.
set -euo pipefail

# shellcheck source=branding.sh
source "$(dirname "${BASH_SOURCE[0]}")/branding.sh"

BACKEND_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TERRAFORM_DIR="${BACKEND_ROOT}/deploy/environments/prod/terraform"
KUBECONFIG_PATH="${TERRAFORM_DIR}/terraform_kubeconfig"
PROJECT="$GCP_PROJECT_ID"
REGION="us-central1"
CLUSTER="prod"

cd "$BACKEND_ROOT"

echo "=== Prerequisites check ==="
command -v terraform >/dev/null || { echo "terraform required (https://developer.hashicorp.com/terraform/install)"; exit 1; }
command -v gcloud >/dev/null || { echo "gcloud required"; exit 1; }

echo "=== Creating terraform_kubeconfig (token-based, matches CI) ==="
# Match CI: get-gke-credentials uses use_auth_provider: false (default), embedding short-lived
# tokens instead of exec-based auth. This avoids needing gke-gcloud-auth-plugin locally.
# Token expires in ~1 hour; re-run the script if a long apply exceeds that.
ENDPOINT="$(gcloud container clusters describe "$CLUSTER" --region "$REGION" --project "$PROJECT" --format='value(endpoint)')"
CA_CERT="$(gcloud container clusters describe "$CLUSTER" --region "$REGION" --project "$PROJECT" --format='value(masterAuth.clusterCaCertificate)' | tr -d '\n')"
TOKEN="$(gcloud auth print-access-token --scopes='https://www.googleapis.com/auth/cloud-platform')"

# Build kubeconfig with embedded token (no exec plugin)
cat > "$KUBECONFIG_PATH" << EOF
apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: ${CA_CERT}
    server: https://${ENDPOINT}
  name: prod_context
contexts:
- context:
    cluster: prod_context
    user: prod_context
    namespace: prod
  name: prod_context
current-context: prod_context
users:
- name: prod_context
  user:
    token: ${TOKEN}
EOF
chmod 600 "$KUBECONFIG_PATH"

echo "=== Terraform init ==="
(cd "$TERRAFORM_DIR" && terraform init -upgrade)

echo "=== Terraform validate ==="
(cd "$TERRAFORM_DIR" && terraform validate)

if [[ "${1:-}" == "-auto-approve" ]]; then
  echo "=== Terraform apply (auto-approve) ==="
  (cd "$TERRAFORM_DIR" && terraform apply -auto-approve "${@:2}")
else
  echo "=== Terraform plan ==="
  (cd "$TERRAFORM_DIR" && terraform plan "$@")
  echo ""
  echo "To apply: $0 -auto-approve [extra terraform args...]"
fi
