variable "cloud_sql_instance_name" {
  default = "postgres-terraform-7"
}

variable "region" {
  type = string
  default = ""
}

variable "postgres_username" {
  default = "postgres"
}

variable "GCP_POSTGRES_PASSWORD" {
  //a TF_VAR_GCP_POSTGRES_PASSWORD environment variable populated value
  type = string  
  default = ""
}

variable "postgres_database_name" {
  default = "notifications"
}

variable "cloud_sql_proxy" {
  default = "/home/joe/cloud_sql_proxy"
}

variable "cloud_sql_proxy_port" {
  default = "1234"
}