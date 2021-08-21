resource "google_artifact_registry_repository" "sockets_repo" {
  provider = google-beta

  project = var.project
  location = var.region
  repository_id = var.repository_id
  description = "go-postgres-sockets Docker repository"
  format = "DOCKER"

  timeouts {
    create = "12m"
  }

  //push a new cloud build for the application into the repository with the cloud_build.sh script
  provisioner "local-exec" {
    command = "${path.module}/cloud_build.sh ${var.region} ${var.project} ${var.repository_id} ${var.image_name}"
  }
}