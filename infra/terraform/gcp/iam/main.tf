resource "google_project_service" "enabled_apis" {
  for_each = toset([
    "iam.googleapis.com",
    "artifactregistry.googleapis.com",
  ])

  project            = var.project_id
  service            = each.value
  disable_on_destroy = true
}
