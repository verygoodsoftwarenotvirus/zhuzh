data "grafana_data_source" "prometheus" {
  name = "grafanacloud-${local.company_slug_ns}-prom"
}

data "grafana_data_source" "loki" {
  name = "grafanacloud-${local.company_slug_ns}-logs"
}

locals {
  prometheus_ds_uid = data.grafana_data_source.prometheus.uid
  loki_ds_uid      = data.grafana_data_source.loki.uid
}

resource "grafana_folder" "alerts" {
  title = "${local.company_name} Alerts"
}

resource "grafana_contact_point" "default" {
  name = "${local.company_name} Default"

  lifecycle {
    create_before_destroy = true
  }

  email {
    addresses               = ["verygoodsoftwarenotvirus@protonmail.com"]
    disable_resolve_message = true
  }
}

resource "grafana_notification_policy" "default" {
  contact_point = grafana_contact_point.default.name
  group_by      = ["alertname"]

  group_wait      = "30s"
  group_interval  = "5m"
  repeat_interval = "4h"
}

# # -----------------------------------------------------------------------------
# # Rule Group 1: Pod Health (PromQL)
# # -----------------------------------------------------------------------------

# resource "grafana_rule_group" "pod_health" {
#   name             = "Pod Health"
#   folder_uid       = grafana_folder.alerts.uid
#   interval_seconds = 60
#   disable_provenance = true

#   rule {
#     name          = "Pod not ready"
#     condition     = "C"
#     for           = "5m"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.prometheus_ds_uid
#       relative_time_range {
#         from = 300
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "kube_pod_status_ready{namespace=\"prod\", condition=\"true\"} == 0"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "Pod {{ $labels.pod }} in namespace {{ $labels.namespace }} is not ready"
#       description = "Pod has been not-ready for more than 5 minutes."
#     }

#     labels = {
#       severity = "warning"
#     }
#   }

#   rule {
#     name          = "Pod crash looping"
#     condition     = "C"
#     for           = "0s"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.prometheus_ds_uid
#       relative_time_range {
#         from = 900
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "increase(kube_pod_container_status_restarts_total{namespace=\"prod\"}[15m]) > 3"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "Container {{ $labels.container }} in pod {{ $labels.pod }} is crash looping"
#       description = "Container has restarted more than 3 times in the last 15 minutes."
#     }

#     labels = {
#       severity = "critical"
#     }
#   }

#   rule {
#     name          = "Pod not running"
#     condition     = "C"
#     for           = "3m"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.prometheus_ds_uid
#       relative_time_range {
#         from = 300
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "kube_pod_status_phase{namespace=\"prod\", phase=~\"Failed|Unknown\"} > 0"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "Pod {{ $labels.pod }} is in {{ $labels.phase }} phase"
#       description = "Pod has been in Failed or Unknown state for more than 3 minutes."
#     }

#     labels = {
#       severity = "critical"
#     }
#   }
# }

# # -----------------------------------------------------------------------------
# # Rule Group 2: Async Message Handler (PromQL + LogQL)
# # -----------------------------------------------------------------------------

# resource "grafana_rule_group" "async_handler" {
#   name             = "Async Message Handler"
#   folder_uid       = grafana_folder.alerts.uid
#   interval_seconds = 60
#   disable_provenance = true

#   rule {
#     name          = "Async handler unavailable"
#     condition     = "C"
#     for           = "3m"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.prometheus_ds_uid
#       relative_time_range {
#         from = 300
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "kube_deployment_status_replicas_available{namespace=\"prod\", deployment=\"${local.company_slug}-async-message-handler-deployment\"} == 0"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "Async message handler has zero available replicas"
#       description = "The async message handler deployment has had no available replicas for more than 3 minutes. Pub/Sub messages are not being consumed."
#     }

#     labels = {
#       severity = "critical"
#     }
#   }

#   rule {
#     name          = "Async handler errors"
#     condition     = "C"
#     for           = "0s"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.loki_ds_uid
#       relative_time_range {
#         from = 300
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "count_over_time({namespace=\"prod\", service_name=~\".*backend.services.*\"} |~ \"error|panic\" [5m])"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "Async message handler is logging errors"
#       description = "The async message handler has logged error or panic messages in the last 5 minutes."
#     }

#     labels = {
#       severity = "warning"
#     }
#   }
# }

# # -----------------------------------------------------------------------------
# # Rule Group 3: CronJob Health (PromQL)
# # -----------------------------------------------------------------------------

# resource "grafana_rule_group" "cronjob_health" {
#   name             = "CronJob Health"
#   folder_uid       = grafana_folder.alerts.uid
#   interval_seconds = 300
#   disable_provenance = true

#   rule {
#     name          = "CronJob not scheduled"
#     condition     = "C"
#     for           = "0s"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.prometheus_ds_uid
#       relative_time_range {
#         from = 86400
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "time() - kube_cronjob_status_last_schedule_time{namespace=\"prod\"} > 86400"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "CronJob {{ $labels.cronjob }} has not been scheduled in 24h"
#       description = "CronJob has not run in the last 24 hours. Check for scheduling issues or suspended state."
#     }

#     labels = {
#       severity = "warning"
#     }
#   }

#   rule {
#     name          = "Job failed"
#     condition     = "C"
#     for           = "5m"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.prometheus_ds_uid
#       relative_time_range {
#         from = 600
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "kube_job_status_failed{namespace=\"prod\"} > 0"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "Job {{ $labels.job_name }} has failed"
#       description = "A Kubernetes Job in the prod namespace has reported failure for more than 5 minutes."
#     }

#     labels = {
#       severity = "warning"
#     }
#   }
# }

# # -----------------------------------------------------------------------------
# # Rule Group 4: Pub/Sub Health (PromQL)
# # -----------------------------------------------------------------------------

# resource "grafana_rule_group" "pubsub_health" {
#   name             = "Pub/Sub Health"
#   folder_uid       = grafana_folder.alerts.uid
#   interval_seconds = 120
#   disable_provenance = true

#   rule {
#     name          = "Dead letter messages"
#     condition     = "C"
#     for           = "5m"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.prometheus_ds_uid
#       relative_time_range {
#         from = 600
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "pubsub_googleapis_com_subscription_num_undelivered_messages{subscription_id=~\".*deadletter.*\"} > 0"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "Dead letter subscription {{ $labels.subscription_id }} has undelivered messages"
#       description = "Messages have landed in a dead letter subscription, indicating processing failures in the primary consumer."
#     }

#     labels = {
#       severity = "critical"
#     }
#   }

#   rule {
#     name          = "Subscription backlog growing"
#     condition     = "C"
#     for           = "10m"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.prometheus_ds_uid
#       relative_time_range {
#         from = 900
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "pubsub_googleapis_com_subscription_num_undelivered_messages{subscription_id!~\".*deadletter.*\"} > 100"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "Subscription {{ $labels.subscription_id }} backlog is growing"
#       description = "Subscription has had more than 100 undelivered messages for over 10 minutes."
#     }

#     labels = {
#       severity = "warning"
#     }
#   }

#   rule {
#     name          = "Old unacked messages"
#     condition     = "C"
#     for           = "5m"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.prometheus_ds_uid
#       relative_time_range {
#         from = 600
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "pubsub_googleapis_com_subscription_oldest_unacked_message_age{subscription_id!~\".*deadletter.*\"} > 3600"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "Subscription {{ $labels.subscription_id }} has messages stuck for over 1 hour"
#       description = "The oldest unacknowledged message is over 1 hour old, indicating the consumer may be stuck or not processing."
#     }

#     labels = {
#       severity = "critical"
#     }
#   }
# }

# # -----------------------------------------------------------------------------
# # Rule Group 5: Cloud SQL Health (PromQL)
# # -----------------------------------------------------------------------------

# resource "grafana_rule_group" "cloudsql_health" {
#   name             = "Cloud SQL Health"
#   folder_uid       = grafana_folder.alerts.uid
#   interval_seconds = 300
#   disable_provenance = true

#   rule {
#     name          = "Disk utilization high"
#     condition     = "C"
#     for           = "10m"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.prometheus_ds_uid
#       relative_time_range {
#         from = 900
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "cloudsql_googleapis_com_database_disk_utilization > 0.85"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "Cloud SQL disk utilization is above 85%"
#       description = "Database disk usage has been above 85% for more than 10 minutes. Consider increasing disk size to avoid running out of space."
#     }

#     labels = {
#       severity = "critical"
#     }
#   }

#   rule {
#     name          = "CPU utilization high"
#     condition     = "C"
#     for           = "15m"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.prometheus_ds_uid
#       relative_time_range {
#         from = 1200
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "cloudsql_googleapis_com_database_cpu_utilization > 0.9"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "Cloud SQL CPU utilization is above 90%"
#       description = "Database CPU has been above 90% for more than 15 minutes. This may indicate query performance issues or an under-provisioned instance."
#     }

#     labels = {
#       severity = "warning"
#     }
#   }

#   rule {
#     name          = "Too many connections"
#     condition     = "C"
#     for           = "5m"
#     no_data_state = "OK"

#     data {
#       ref_id         = "A"
#       datasource_uid = local.prometheus_ds_uid
#       relative_time_range {
#         from = 600
#         to   = 0
#       }
#       model = jsonencode({
#         refId         = "A"
#         expr          = "cloudsql_googleapis_com_database_postgresql_num_backends > 80"
#         intervalMs    = 1000
#         maxDataPoints = 43200
#       })
#     }

#     data {
#       ref_id         = "C"
#       datasource_uid = "-100"
#       relative_time_range {
#         from = 0
#         to   = 0
#       }
#       model = jsonencode({
#         refId      = "C"
#         type       = "threshold"
#         expression = "A"
#         conditions = [{
#           evaluator = {
#             type   = "gt"
#             params = [0]
#           }
#         }]
#       })
#     }

#     annotations = {
#       summary     = "Cloud SQL connection count is approaching the limit"
#       description = "PostgreSQL backend connections exceed 80 (default max is 100). Risk of connection exhaustion."
#     }

#     labels = {
#       severity = "warning"
#     }
#   }
# }

# -----------------------------------------------------------------------------
# Rule Group 6: Error Rates (LogQL)
# -----------------------------------------------------------------------------

resource "grafana_rule_group" "error_rates" {
  name             = "Error Rates"
  folder_uid       = grafana_folder.alerts.uid
  interval_seconds = 60
  disable_provenance = true

  rule {
    name          = "High error rate"
    condition     = "C"
    for           = "0s"
    no_data_state = "OK"

    data {
      ref_id         = "A"
      datasource_uid = local.loki_ds_uid
      relative_time_range {
        from = 300
        to   = 0
      }
      model = jsonencode({
        refId         = "A"
        expr          = "sum(count_over_time({namespace=\"prod\"} |= \"level\" |= \"error\" [5m])) > 50"
        intervalMs    = 1000
        maxDataPoints = 43200
      })
    }

    data {
      ref_id         = "C"
      datasource_uid = "-100"
      relative_time_range {
        from = 0
        to   = 0
      }
      model = jsonencode({
        refId      = "C"
        type       = "threshold"
        expression = "A"
        conditions = [{
          evaluator = {
            type   = "gt"
            params = [0]
          }
        }]
      })
    }

    annotations = {
      summary     = "High error log volume in prod namespace"
      description = "More than 50 error-level log entries have been recorded in the last 5 minutes across all prod services."
    }

    labels = {
      severity = "warning"
    }
  }

  rule {
    name          = "Database errors"
    condition     = "C"
    for           = "0s"
    no_data_state = "OK"

    data {
      ref_id         = "A"
      datasource_uid = local.loki_ds_uid
      relative_time_range {
        from = 300
        to   = 0
      }
      model = jsonencode({
        refId         = "A"
        expr          = "count_over_time({namespace=\"prod\"} |~ \"database|connection refused|sql\" [5m]) > 0"
        intervalMs    = 1000
        maxDataPoints = 43200
      })
    }

    data {
      ref_id         = "C"
      datasource_uid = "-100"
      relative_time_range {
        from = 0
        to   = 0
      }
      model = jsonencode({
        refId      = "C"
        type       = "threshold"
        expression = "A"
        conditions = [{
          evaluator = {
            type   = "gt"
            params = [0]
          }
        }]
      })
    }

    annotations = {
      summary     = "Database connectivity errors detected"
      description = "Logs containing database-related error patterns (database, connection refused, sql) have been detected in the last 5 minutes."
    }

    labels = {
      severity = "critical"
    }
  }
}
