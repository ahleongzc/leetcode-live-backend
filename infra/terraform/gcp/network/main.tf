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
}
