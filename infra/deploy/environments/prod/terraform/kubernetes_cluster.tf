# Secret Manager API must be enabled for the CSI driver and SecretSync to fetch secrets
resource "google_project_service" "secret_manager" {
  project = local.gcp_project_id
  service = "secretmanager.googleapis.com"
}

resource "google_project_iam_custom_role" "prod_cluster_role" {
  role_id     = "prod_cluster_role"
  title       = "Prod cluster role"
  description = "An IAM role for the prod cluster"
  permissions = [
    "secretmanager.versions.access",
    "cloudsql.instances.connect",
    "cloudsql.instances.get",
    "pubsub.topics.list",
    "pubsub.topics.publish",
    "cloudtrace.traces.patch",
    "logging.buckets.create",
    "logging.buckets.write",
    "logging.buckets.list",
    "logging.buckets.get",
    "storage.objects.list",
    "storage.objects.get",
    "storage.objects.update",
    "storage.objects.create",
    "storage.objects.delete",
    "storage.objects.get",
  ]
}

resource "google_service_account" "prod_cluster_service_account" {
  account_id   = "prod-cluster"
  display_name = "prod Cluster Service Account"
}

resource "google_project_iam_member" "prod_cluster" {
  project = local.gcp_project_id
  role    = google_project_iam_custom_role.prod_cluster_role.id
  member  = format("serviceAccount:%s", google_service_account.prod_cluster_service_account.email)
}

resource "google_container_cluster" "primary" {
  provider = google-beta
  name     = local.environment
  location = local.gcp_region

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1
  deletion_protection      = false

  network = google_compute_network.private_network.name
  # subnetwork = google_compute_subnetwork.subnet.name

  # Enable the Secret Manager add-on for native GCP Secret Manager integration
  # This allows pods to mount secrets from GCP Secret Manager as volumes
  secret_manager_config {
    enabled = true
  }

  # Enable SecretSync to sync GCP secrets into native Kubernetes Secrets
  secret_sync_config {
    enabled = true
  }

  # Workload Identity is required for the Secret Manager add-on
  workload_identity_config {
    workload_pool = "${local.gcp_project_id}.svc.id.goog"
  }
}

# Separately Managed Node Pool
resource "google_container_node_pool" "primary_nodes" {
  location   = local.gcp_region
  name       = google_container_cluster.primary.name
  cluster    = google_container_cluster.primary.name
  node_count = 1

  node_locations = [
    local.gcp_main_zone
  ]

  autoscaling {
    total_min_node_count = 1
    total_max_node_count = 1
    location_policy = "ANY"
  }

  management {
    auto_repair = true
    auto_upgrade = true
  }

  network_config {
    
  }

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]

    labels = {
      env = local.gcp_project_id
    }

    preemptible  = false
    machine_type = "e2-standard-2"
    tags         = ["gke-node", local.environment]
    metadata = {
      disable-legacy-endpoints = "true"
    }

    # Use Workload Identity for pod-level GCP authentication
    workload_metadata_config {
      mode = "GKE_METADATA"
    }
  }
}

# Kubernetes service account for workloads that need to access GCP Secret Manager
resource "google_service_account" "workload_identity_sa" {
  account_id   = "workload-identity-sa"
  display_name = "Workload Identity Service Account"
  description  = "Service account for Kubernetes workloads to access GCP resources"
}

# Grant the workload identity SA access to secrets and Pub/Sub
resource "google_project_iam_member" "workload_identity_secret_accessor" {
  project = local.gcp_project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.workload_identity_sa.email}"
}

resource "google_project_iam_member" "workload_identity_pubsub" {
  project = local.gcp_project_id
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${google_service_account.workload_identity_sa.email}"
}

resource "google_project_iam_member" "workload_identity_pubsub_subscriber" {
  project = local.gcp_project_id
  role    = "roles/pubsub.subscriber"
  member  = "serviceAccount:${google_service_account.workload_identity_sa.email}"
}

# TopicAdminClient.GetTopic requires pubsub.topics.get (included in viewer)
resource "google_project_iam_member" "workload_identity_pubsub_viewer" {
  project = local.gcp_project_id
  role    = "roles/pubsub.viewer"
  member  = "serviceAccount:${google_service_account.workload_identity_sa.email}"
}

resource "google_project_iam_member" "workload_identity_monitoring_viewer" {
  project = local.gcp_project_id
  role    = "roles/monitoring.viewer"
  member  = "serviceAccount:${google_service_account.workload_identity_sa.email}"
}

# Firebase Cloud Messaging for Android push notifications (async message handler)
resource "google_project_iam_member" "workload_identity_fcm_admin" {
  project = local.gcp_project_id
  role    = "roles/firebasecloudmessaging.admin"
  member  = "serviceAccount:${google_service_account.workload_identity_sa.email}"
}

# Allow the Kubernetes service account to impersonate the GCP service account
resource "google_service_account_iam_member" "workload_identity_binding" {
  service_account_id = google_service_account.workload_identity_sa.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${local.gcp_project_id}.svc.id.goog[prod/api-service-account]"
}
