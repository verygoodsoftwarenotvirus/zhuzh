data "google_iam_policy" "public_policy" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
      "allUsers",
    ]
  }
}

# Media bucket policy: public read + workload-identity-sa write (for API uploads)
data "google_iam_policy" "api_media_policy" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
      "allUsers",
    ]
  }
  binding {
    role = "roles/storage.objectAdmin"
    members = [
      "serviceAccount:workload-identity-sa@${local.gcp_project_id}.iam.gserviceaccount.com",
    ]
  }
}
