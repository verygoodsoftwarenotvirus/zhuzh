variable "GRAFANA_AUTH_TOKEN" {}
variable "GRAFANA_URL" {}

provider "grafana" {
  url  = var.GRAFANA_URL
  auth = var.GRAFANA_AUTH_TOKEN
}
