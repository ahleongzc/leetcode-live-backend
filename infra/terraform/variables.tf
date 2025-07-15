variable "cloudflare_api_token" {
  description = "cloudflare api token"
  type = string
  sensitive = true
}

variable "gcp_sa_key_file_path" {
  description = "service acount key file path"
  type        = string
}

variable "db_password" {
  description = "database password"
  type        = string
  sensitive   = true
}

variable "project_id" {
  description = "gcp project id"
  type        = string
}

variable "region" {
  description = "gcp region"
  type        = string
}

variable "zone" {
  description = "gcp zone"
  type        = string
}

variable "physical_database_name" {
  description = "physical database name"
  type        = string
}

variable "database_name" {
  description = "database name"
  type        = string
}

variable "database_user" {
  description = "database user"
  type        = string
}

variable "database_version" {
  description = "database version"
  type        = string
}

variable "whitelisted_ip_address" {
  description = "whitelisted ip address"
  type        = string
  sensitive   = true
}
