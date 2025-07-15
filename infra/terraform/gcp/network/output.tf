output "vpc_id" {
  value = google_compute_network.vpc.id
}
output "private_subnet_id" {
  value = google_compute_subnetwork.private_subnet.id
}
