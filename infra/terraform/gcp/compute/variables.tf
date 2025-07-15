variable "bastion_name" {
  default = "bastion-host"
}
variable "bastion_machine_type" {
  default = "e2-micro"
}
// TODO: Change this to a more powerful instance
variable "application_server_machine_type" {
  default = "e2-micro"
}
variable "bastion_image_family" {
  default = "debian-12"
}
variable "region" {
  type = string
}
variable "zone" {
  type = string
}
variable "public_subnet_self_link" {
  description = "the self link of the public subnet where the bastion host will reside."
  type        = string
}
variable "ssh_public_key" {
  type = string
}
