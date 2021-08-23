data "google_client_config" "default" {}

provider "kubernetes" {
  host                   = "https://${module.gke.endpoint}"
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(module.gke.ca_certificate)
}

module "gke" {
  source                   = "terraform-google-modules/kubernetes-engine/google"
  project_id               = var.project_id
  name                     = "terraform-gke"
  region                   = var.region
  network                  = var.network
  subnetwork               = var.subnetwork
  ip_range_pods            = var.ip_range_pods
  ip_range_services        = var.ip_range_services
  remove_default_node_pool = true
  service_account          = "create"
  node_metadata            = "GKE_METADATA_SERVER"
  node_pools = [
    {
      name         = "default-pool"
      min_count    = 1
      max_count    = 1
      auto_upgrade = false
      image_type   = "COS_CONTAINERD"
    }
  ]
}

module "workload_identity" {
  source              = "terraform-google-modules/kubernetes-engine/google//modules/workload-identity"
  project_id          = var.project_id
  name                = "workload-identity-${module.gke.name}"
  namespace           = "default"
  roles               = ["roles/cloudsql.client"]
}