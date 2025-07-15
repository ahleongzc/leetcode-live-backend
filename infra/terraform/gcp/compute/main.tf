resource "google_compute_instance" "bastion" {
  name         = var.bastion_name
  machine_type = var.bastion_machine_type
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = data.google_compute_image.debian.self_link
    }
  }

  network_interface {
    subnetwork = var.public_subnet_self_link
    access_config {}
  }

  metadata = {
    ssh-keys = "dev:${var.ssh_public_key}"
  }
  tags = ["bastion-host"]
  labels = {
    role = "bastion"
  }
}

resource "google_compute_instance" "application_server" {
  name         = "application-server"
  machine_type = var.application_server_machine_type
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = data.google_compute_image.debian.self_link
    }
  }

  network_interface {
    subnetwork = var.public_subnet_self_link
    access_config {}
  }

  metadata = {
    ssh-keys = "dev:${var.ssh_public_key}"
  }

  tags = ["application-server"]

  labels = {
    role = "application"
  }
}


data "google_compute_image" "debian" {
  family  = var.bastion_image_family
  project = "debian-cloud"
}
