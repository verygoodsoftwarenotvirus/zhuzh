locals {
  company_name         = "Zhuzh"
  company_slug         = "zhuzh"
  company_slug_ns      = "zhuzh"
  public_domain        = "${local.company_slug_ns}.dev"
  admin_domain         = "admin.${local.public_domain}"
  media_domain         = "media.${local.public_domain}"
  userdata_domain      = "userdata.${local.public_domain}"
  gcp_project_id       = "${local.company_slug}-prod"
  tf_cloud_org         = local.company_slug_ns
  ios_bundle_id        = "com.${local.company_slug_ns}.ios"
  k8s_admin_webapp_cfg = "${local.company_slug}-admin-webapp-config"
  k8s_consumer_webapp  = "${local.company_slug}-consumer-webapp-secrets"
}
