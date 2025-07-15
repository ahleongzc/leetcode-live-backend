// TODO
// 1. Set up SSL configuration
// 2. Set up backups

resource "google_compute_global_address" "private_ip_address" {
  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = var.vpc_id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = var.vpc_id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

resource "google_sql_database_instance" "database_server" {
  name                = var.database_name
  database_version    = var.database_version
  region              = var.region
  deletion_protection = false

  settings {
    tier = "db-f1-micro"

    ip_configuration {
      ipv4_enabled    = false
      private_network = var.vpc_id
    }

    backup_configuration {
      enabled = false
    }

    maintenance_window {
      day          = 7
      hour         = 3
      update_track = "stable"
    }

    availability_type = "ZONAL"
  }

  depends_on = [google_service_networking_connection.private_vpc_connection]
}

resource "google_sql_database" "physical_database" {
  name     = var.physical_database_name
  instance = google_sql_database_instance.database_server.id

  depends_on = [google_sql_database_instance.database_server]
}

resource "google_sql_user" "database_user" {
  name     = var.database_user
  password = var.db_password
  instance = google_sql_database_instance.database_server.id

  depends_on = [google_sql_database_instance.database_server]
}
