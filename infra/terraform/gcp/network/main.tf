resource "google_compute_network" "vpc" {
  name                    = "prod"
  auto_create_subnetworks = false
  routing_mode            = "GLOBAL"
  project                 = var.project_id
}

resource "google_compute_subnetwork" "public_subnet" {
  name                     = "public-subnet"
  project                  = var.project_id
  ip_cidr_range            = var.public_subnet_cidr
  region                   = var.region
  network                  = google_compute_network.vpc.id
  description              = "Public subnet with external IPs allowed"
  private_ip_google_access = true
}

resource "google_compute_subnetwork" "private_subnet" {
  name                     = "private-subnet"
  project                  = var.project_id
  ip_cidr_range            = var.private_subnet_cidr
  region                   = var.region
  network                  = google_compute_network.vpc.id
  description              = "Private subnet for internal instances without external IPs"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "google-managed-services"
    ip_cidr_range = var.secondary_private_subnet_cidr
  }
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network = google_compute_network.vpc.id
  service = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [
    google_compute_subnetwork.private_subnet.secondary_ip_range[0].range_name
  ]

  depends_on = [
    google_compute_subnetwork.private_subnet
  ]
}
