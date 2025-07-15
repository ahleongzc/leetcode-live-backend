resource "google_sql_database_instance" "prod" {
  name             = var.database_name
  database_version = var.database_version
  region           = var.region

  settings {
    tier = "db-f1-micro"

    ip_configuration {
      ipv4_enabled = true

      authorized_networks {
        value = var.whitelisted_ip_address
      }
    }

    backup_configuration {
      enabled = false
    }

    availability_type = "ZONAL"
  }

  deletion_protection = false
}

resource "google_sql_database" "database" {
  name     = var.physical_database_name
  instance = google_sql_database_instance.prod.name
}

resource "google_sql_user" "database_user" {
  name     = var.database_user
  instance = google_sql_database_instance.prod.name
  password = var.db_password
}
