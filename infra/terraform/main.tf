module "gcp_sql" {
  source                 = "./gcp/sql"
  region                 = var.region
  database_user          = var.database_user
  db_password            = var.db_password
  physical_database_name = var.physical_database_name
  database_version       = var.database_version
  whitelisted_ip_address = var.whitelisted_ip_address
  database_name          = var.database_name
}

module "gcp_backend_instance" {
  source = "./gcp/compute"
  zone   = var.zone
}

module "gcp_network" {
  source              = "./gcp/network"
  project_id          = var.project_id
  region              = var.region
  public_subnet_cidr  = var.public_subnet_cidr
  private_subnet_cidr = var.private_subnet_cidr

}
