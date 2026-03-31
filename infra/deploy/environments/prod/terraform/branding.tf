locals {
  company_name     = "Zhuzh"
  company_slug     = "zhuzh"
  company_slug_ns  = "zhuzh"
  public_domain    = "${local.company_slug_ns}.dev"
  gcp_project_id   = "${local.company_slug}-prod"
}
