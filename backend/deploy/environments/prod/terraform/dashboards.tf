resource "grafana_folder" "dashboards" {
  title = "${local.company_name} Dashboards"
}

resource "grafana_dashboard" "api_server" {
  folder = grafana_folder.dashboards.id

  config_json = file("${path.module}/dashboards/api_server.json")
}

resource "grafana_dashboard" "admin_webapp" {
  folder = grafana_folder.dashboards.id

  config_json = file("${path.module}/dashboards/admin_webapp.json")
}

resource "grafana_dashboard" "consumer_webapp" {
  folder = grafana_folder.dashboards.id

  config_json = file("${path.module}/dashboards/consumer_webapp.json")
}

resource "grafana_dashboard" "async_message_handler" {
  folder = grafana_folder.dashboards.id

  config_json = file("${path.module}/dashboards/async_message_handler.json")
}

resource "grafana_dashboard" "gcp_infrastructure" {
  folder = grafana_folder.dashboards.id

  config_json = file("${path.module}/dashboards/gcp_infrastructure.json")
}
