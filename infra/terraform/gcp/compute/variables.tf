variable "backend_instance_name" {
  description = "name of the instance"
  type        = string
  default     = "backend"
}

variable "machine_type" {
  description = "gce machine type"
  type        = string
  default     = "e2-micro"
}

variable "zone" {
  description = "gcp zone"
  type        = string
}
