resource "google_artifact_registry_repository" "container_registry" {
  provider      = google
  location      = var.location
  repository_id = var.repository_id
  description   = "container registry for docker images"
  format        = "DOCKER"
  mode          = "STANDARD_REPOSITORY"

  docker_config {
    immutable_tags = true
  }

  depends_on = [var.enabled_apis]
}

resource "google_service_account" "cicd_service_account" {
  account_id   = "github-actions"
  display_name = "service account for github actions"
  project      = var.project_id

  depends_on = [var.enabled_apis]
}

resource "google_artifact_registry_repository_iam_member" "artifact_registry_writer_binding" {
  project    = var.project_id
  role       = "roles/artifactregistry.writer"
  repository = google_artifact_registry_repository.container_registry.repository_id
  member     = "serviceAccount:${google_service_account.cicd_service_account.email}"

  depends_on = [google_artifact_registry_repository.container_registry]
}
