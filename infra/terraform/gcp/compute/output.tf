output "application_server_ip" {
  value       = google_compute_instance.application_server.network_interface[0].access_config[0].nat_ip
  description = "The external IP address of the application server"
}
