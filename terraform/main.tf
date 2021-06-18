provider "google" {
  credentials = "${file("service-account-key.json")}"
  project     = "${var.gcp_project}"
  region      = "${var.gcp_region}"
}

provider "google-beta" {
  credentials = "${file("service-account-key.json")}"
  project     = "${var.gcp_project}"
  region      = "${var.gcp_region}"
}

resource "google_compute_network" "private_network" {
  provider     = "google-beta"
  name         = "private-network"
  routing_mode = "REGIONAL"
}

resource "google_compute_global_address" "private_ip_address" {
  provider      = "google-beta"
  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.private_network.self_link
}

resource "google_service_networking_connection" "private_vpc_connection" {
  provider = google-beta
  network                 = google_compute_network.private_network.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

/*Investigar este betulio*/
resource "google_vpc_access_connector" "vpc_serverless_connector" {
  name          = "serverless_vpc"
  ip_cidr_range = "10.8.0.0/28"
  network       = google_compute_network.private_network.self_link
}


# Instancia de cloud sql
resource "google_sql_database_instance" "cloud_sql" {
  provider    = google-beta
  name        = "${var.database_name}"
  database_version = "POSTGRES_11"
  depends_on  = [google_service_networking_connection.private_vpc_connection]
  settings {
    tier = "db-f1-micro"
    user_labels = {
      Name = "${var.database_name}"
      Environment = "${var.environment}"
      Tier = "database"
      Type = "Postgres11"
    }
    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.private_network.self_link
    }
  }
}
