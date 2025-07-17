output "application_server_ip" {
  value       = module.gcp_compute.application_server_ip
  description = "The external IP address of the application server"
}

output "ssh_user" {
  value       = var.ssh_user
  description = "SSH user"
}

output "ssh_public_key_path" {
  value = var.ssh_public_key_path
}