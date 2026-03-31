#!/usr/bin/env bash
# Run terraform plan -destroy/destroy for prod infra (GKE cluster, networking, etc.).
# Prerequisites: terraform, gcloud, Terraform Cloud login (terraform login), workspace variables set.
set -euo pipefail

INFRA_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TERRAFORM_DIR="${INFRA_ROOT}/deploy/environments/prod/terraform"

cd "$INFRA_ROOT"

echo "=== Prerequisites check ==="
command -v terraform >/dev/null || { echo "terraform required (https://developer.hashicorp.com/terraform/install)"; exit 1; }
command -v gcloud >/dev/null || { echo "gcloud required"; exit 1; }

echo "=== Terraform init ==="
(cd "$TERRAFORM_DIR" && terraform init -upgrade)

echo "=== Terraform validate ==="
(cd "$TERRAFORM_DIR" && terraform validate)

if [[ "${1:-}" == "-auto-approve" ]]; then
  echo "=== Terraform destroy (auto-approve) ==="
  (cd "$TERRAFORM_DIR" && terraform destroy -auto-approve "${@:2}")
else
  echo "=== Terraform plan -destroy ==="
  (cd "$TERRAFORM_DIR" && terraform plan -destroy "$@")
  echo ""
  echo "To destroy: $0 -auto-approve [extra terraform args...]"
fi
