variable "project" {
  type = string
  default = ""
}

variable "region" {
  type = string
  default = ""
}

variable "repository_id" {
  default = "go-postgres-sockets-repo"
}

variable "application_name" {
  default = "go-postgres-sockets-terraform"
}

variable "image_name" {
  default = "sockets-app"
}