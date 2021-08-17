variable "project" { 
  default = "bebo-dev-sockets-terraform"
}

variable "gcp_credentials_file" {
  default = "./gcp-sa-credentials.json"
}

variable "region" {
  default = "europe-west2"
}

variable "zone" {
  default = "europe-west2-b"
}

variable "GCP_POSTGRES_PASSWORD" {
  //a TF_VAR_GCP_POSTGRES_PASSWORD environment variable populated value
  type = string  
  default = ""
}