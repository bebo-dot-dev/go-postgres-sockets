terraform {
  required_providers {    
    postgresql = {
      source = "cyrilgdn/postgresql"
      version = "1.13.0"
    }
  }
}

provider "postgresql" {
  host            = var.postgres_server_ip
  port            = 5432
  username        = var.postgres_username
  password        = var.GCP_POSTGRES_PASSWORD
  sslmode         = "require"
}

resource "google_sql_database_instance" "instance" {
  name             = var.cloud_sql_instance_name
  database_version = "POSTGRES_13"
  region           = var.region

  settings {
    tier = "db-f1-micro"
  }  
}

resource "google_sql_database" "database" {
  name     = var.postgres_database_name
  instance = google_sql_database_instance.instance.name
}

resource "google_sql_user" "users" {
  name     = var.postgres_username
  instance = google_sql_database_instance.instance.name
  password = var.GCP_POSTGRES_PASSWORD

  //restore the db with the restore_db.sh script via psql and cloud_sql_proxy
  provisioner "local-exec" {
    command = "${path.module}/restore_db.sh ${var.cloud_sql_proxy} ${google_sql_database_instance.instance.connection_name} ${var.cloud_sql_proxy_port} ${var.postgres_username} ${var.GCP_POSTGRES_PASSWORD}"
  }
}