variable "region" {
  description = "region for sql instance"
  type        = string
}

variable "database_user" {
  description = "database user"
  type        = string
}

variable "db_password" {
  description = "database password"
  type        = string
  sensitive   = true
}

variable "physical_database_name" {
  description = "database name"
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

variable "database_name" {
  description = "database name"
  type        = string
}
