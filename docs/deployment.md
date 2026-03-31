# Deployment Guide

## Release-based deployment

- **Default branch:** `main` — clone and work from `main`. All PRs and CI target `main`.
- **Production deploy trigger:** **Publishing a GitHub Release** (with a semver tag, e.g. `v1.2.3`) runs the deploy pipeline. The workflow checks out the tag, so the build and deploy use the exact commit that was released.

To deploy to production:

1. Ensure the code you want to ship is merged into `main`.
2. Create a git tag (e.g. `v1.2.3`) on the commit you want to deploy.
3. Create a GitHub Release from that tag and **publish** it. Publishing the release triggers [`.github/workflows/deploy_prod.yaml`](.github/workflows/deploy_prod.yaml).

**Where the version appears:** The release tag is embedded in backend binaries (exposed via version/health endpoints and JSON), included in Grafana deployment and Terraform annotations, and set on Kubernetes resources as the `app.kubernetes.io/version` label.

### Repository settings

- In **Settings → General**, set the default branch to `main`. Create `main` from your current default (e.g. rename `dev` or create from it) and update any open PRs.
- Remove or relax branch protection on `prod` if it exists; you can delete the `prod` branch after cutting over to release-based deploys.

---

## What happens on deploy

Publishing a GitHub Release triggers [`.github/workflows/deploy_prod.yaml`](.github/workflows/deploy_prod.yaml). The pipeline runs three jobs in sequence:

| Job                | Description                                                                                           |
|--------------------|-------------------------------------------------------------------------------------------------------|
| **baseline-infra** | Terraform for GCP baseline (project, GKE cluster, networking, DNS, email, IAM)                        |
| **backend-infra**  | Terraform for backend resources (Cloud SQL, storage buckets, Pub/Sub, Algolia, Kubernetes namespaces) |
| **applications**   | Builds and deploys app workloads via Skaffold (API server, workers, cronjobs)                         |

Infrastructure is split into:

- **`infra/deploy/environments/prod/`** — GKE cluster, networking (Cloudflare), email (SendGrid), Caddy (TLS + reverse proxy), external-dns
- **`backend/deploy/environments/prod/`** — Database, storage, Pub/Sub, search, and application Terraform/Kustomize

---

## Local production deploy (applications only)

To run the application deploy locally (e.g., for debugging or manual rollouts), use:

```bash
./scripts/deploy-prod-local.sh
```

**Prerequisites**: Docker, gcloud, kubectl, skaffold. Your `kubectl` context must point at the prod GKE cluster.

This script only runs the Skaffold application deploy. It does **not** run Terraform or update infrastructure. Version is embedded in binaries via ldflags (set `VERSION` to override the default `local`).

---

## Post-deployment verification

After each deploy, run verification to ensure the deployment is useful:

```bash
# Wait 1-2 min for load balancer propagation, then:
make verify_prod
```

Or manually: `skaffold verify --filename=skaffold.yaml --profile prod` followed by `./scripts/post-deploy-verify.sh`.

The CI pipeline runs these verification steps automatically after `skaffold run`. See [post-deployment-checklist.md](post-deployment-checklist.md) for manual checks (observability, queues, etc.) if that file exists.

---

## Configuration and secrets

Backend configuration and secrets are documented in [backend/docs/configuration.md](backend/docs/configuration.md). Terraform uses Terraform Cloud for variables; application secrets are synced from GCP Secret Manager into the cluster.
