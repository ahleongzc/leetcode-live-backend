// TODO
// 1. Set up SSL configuration
// 2. Set up backups

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
}

resource "google_sql_database" "physical_database" {
  name     = var.physical_database_name
  instance = google_sql_database_instance.database_server.id
}

resource "google_sql_user" "database_user" {
  name     = var.database_user
  password = var.db_password
  instance = google_sql_database_instance.database_server.id
}
