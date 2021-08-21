terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "3.5.0"
    }    
  }
}

provider "google" {
  project = var.project
  credentials = file(var.gcp_credentials_file)
  region  = var.region
  zone    = var.zone
}

module "postgres_database" {
  source = "./modules/postgres-database"
  region = "${var.region}"
  GCP_POSTGRES_PASSWORD = var.GCP_POSTGRES_PASSWORD
}

module "docker_build" {
  source = "./modules/docker-build"
  project = "${var.project}"
  region = "${var.region}"
}