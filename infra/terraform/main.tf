module "gcp_network" {
  source                        = "./gcp/network"
  project_id                    = var.project_id
  region                        = var.region
  public_subnet_cidr            = var.public_subnet_cidr
  private_subnet_cidr           = var.private_subnet_cidr
  secondary_private_subnet_cidr = var.secondary_private_subnet_cidr
}
module "gcp_compute" {
  source                  = "./gcp/compute"
  zone                    = var.zone
  region                  = var.region
  public_subnet_self_link = module.gcp_network.public_subnet_self_link
  ssh_public_key          = file(pathexpand(var.ssh_public_key_path))
  ssh_user                = var.ssh_user
}
module "gcp_sql" {
  source                 = "./gcp/sql"
  region                 = var.region
  database_user          = var.database_user
  db_password            = var.db_password
  physical_database_name = var.physical_database_name
  database_version       = var.database_version
  database_name          = var.database_name
  vpc_id                 = module.gcp_network.vpc_id

  depends_on = [module.gcp_network]
}
