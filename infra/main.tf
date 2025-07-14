module "sql" {
  source                 = "./modules/sql"
  region                 = var.region
  database_user          = var.database_user
  db_password            = var.db_password
  physical_database_name = var.physical_database_name
  database_version       = var.database_version
  whitelisted_ip_address = var.whitelisted_ip_address
  database_name          = var.database_name
}
