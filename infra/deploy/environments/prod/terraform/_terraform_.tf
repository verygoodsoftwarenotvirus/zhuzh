terraform {
  required_version = "1.14.5"

  backend "remote" {
    organization = "zhuzh" # must be literal; backend blocks don't support interpolation

    workspaces {
      name = "prod-infra"
    }
  }
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "5.18.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "7.23.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "7.23.0"
    }
    grafana = {
      source  = "grafana/grafana"
      version = "4.28.0"
    }
  }
}
