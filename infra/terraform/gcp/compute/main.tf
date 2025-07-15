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
  }

  metadata = {
    ssh-keys = "developer:${var.ssh_public_key}"
  }
  tags = ["bastion-host"]
  labels = {
    role = "bastion"
  }
}