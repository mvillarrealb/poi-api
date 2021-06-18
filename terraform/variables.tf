variable "gcp_region" {
  description = "GCP region"
  default = "us-east1"
}

variable "environment" {
  default = "demo"
}

variable "gcp_zone" {
  description = "GCP zone"
  default = "us-east1-a"
}

variable "gcp_project" {
  description = "GCP project name"
  default = "mvillarreal-demo-platform"
}

variable "database_name" {
    default = "mvillarreal-pg-sql"
}
