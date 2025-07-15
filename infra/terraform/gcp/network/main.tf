// TODO
// 1. Add source_ranges to my firewall rules to SSH into the bastion host

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

resource "google_compute_firewall" "allow_ssh_to_bastion" {
  name        = "allow-ssh-to-bastion"
  network     = google_compute_network.vpc.id
  direction   = "INGRESS"
  target_tags = ["bastion-host"]

  # TODO: Change this to home network
  source_ranges = ["0.0.0.0/0"]
  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
}

resource "google_compute_firewall" "allow_bastion_to_cloudsql" {
  name      = "allow-bastion-to-cloudsql"
  network   = google_compute_network.vpc.id
  direction = "EGRESS"
  allow {
    protocol = "tcp"
    ports    = ["5432"]
  }
  destination_ranges = [var.secondary_private_subnet_cidr] # Assuming this is where Cloud SQL lives
}
