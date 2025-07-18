output "artifact_registry_repository_url" {
  value = google_artifact_registry_repository.container_registry
}

output "cicd_service_account_email" {
  value = google_service_account.cicd_service_account.email
}
