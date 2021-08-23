variable "project_id" {}

variable "region" {}

variable "network" {
  description = "The VPC network to host the cluster in"
  default = "default"
}

variable "subnetwork" {
  description = "The subnetwork to host the cluster in"
  default = "default"
}

variable "ip_range_pods" {
  description = "The secondary ip range to use for pods"
  default = ""
}

variable "ip_range_services" {
  description = "The secondary ip range to use for pods"
  default = ""
}