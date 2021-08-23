# go-postgres-sockets

### This directory contains the terraform configuration equivalent of what is implemented with native gcloud and kubectl cli commands in the [k8s](https://github.com/bebo-dot-dev/go-postgres-sockets/tree/main/cmd/go-postgres-sockets/k8s) directory. Applying the terraform configuration in this directory will result in this application being built and hosted in a Google Cloud Platform (GCP) k8s Google Kubernetes Engine (GKE) cluster at https://sockets.bebo.dev

### The terraform configuration in this directory assumes local installs of `glcoud` (GCP SDK) to augment terraform with various direct gcloud commands and `cloud_sql_proxy` and `psql` to enable a restore of the [postgreg db](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/postgres/notifications_db_backup.sql)
---
### 1. create a new GCP project and set it as the default project. The GCP project Id needs to be unique across all projects ever created by anyone in GCP
```
gcloud projects create bebo-dev-sockets-terraform-2 --name="go postgres sockets terraform"

gcloud config set project bebo-dev-sockets-terraform-2
```
### 2. enable billing on the new project. The billing accountId can be found in the GCP billing console
```
gcloud beta billing projects link bebo-dev-sockets-terraform-2 --billing-account=xxxxxx-xxxxxx-xxxxxx
```
### 3. enable required APIs on the new project
```
gcloud services enable artifactregistry.googleapis.com &&
gcloud services enable cloudbuild.googleapis.com &&
gcloud services enable compute.googleapis.com &&
gcloud services enable container.googleapis.com &&
gcloud services enable dataproc.googleapis.com &&
gcloud services enable iamcredentials.googleapis.com &&
gcloud services enable pubsub.googleapis.com &&
gcloud services enable servicenetworking.googleapis.com &&
gcloud services enable sourcerepo.googleapis.com &&
gcloud services enable sqladmin.googleapis.com &&
gcloud services enable cloudresourcemanager.googleapis.com
```
### 4. create a new named Google service account and assign required roles on the new service account
```
gcloud iam service-accounts create terraform-sockets-sa

gcloud projects add-iam-policy-binding bebo-dev-sockets-terraform \
  --member serviceAccount:terraform-sockets-sa@bebo-dev-sockets-terraform.iam.gserviceaccount.com \
  --role roles/editor

gcloud projects add-iam-policy-binding bebo-dev-sockets-terraform \
  --member serviceAccount:terraform-sockets-sa@bebo-dev-sockets-terraform.iam.gserviceaccount.com \
  --role roles/cloudsql.admin

gcloud projects add-iam-policy-binding bebo-dev-sockets-terraform \
  --member serviceAccount:terraform-sockets-sa@bebo-dev-sockets-terraform.iam.gserviceaccount.com \
  --role roles/cloudbuild.builds.editor

gcloud projects add-iam-policy-binding bebo-dev-sockets-terraform \
  --member serviceAccount:terraform-sockets-sa@bebo-dev-sockets-terraform.iam.gserviceaccount.com \
  --role roles/cloudbuild.builds.viewer

gcloud projects add-iam-policy-binding bebo-dev-sockets-terraform \
  --member serviceAccount:terraform-sockets-sa@bebo-dev-sockets-terraform.iam.gserviceaccount.com \
  --role roles/artifactregistry.admin

gcloud projects add-iam-policy-binding bebo-dev-sockets-terraform \
  --member serviceAccount:terraform-sockets-sa@bebo-dev-sockets-terraform.iam.gserviceaccount.com \
  --role roles/resourcemanager.projectIamAdmin

gcloud projects add-iam-policy-binding bebo-dev-sockets-terraform \
  --member serviceAccount:terraform-sockets-sa@bebo-dev-sockets-terraform.iam.gserviceaccount.com \
  --role roles/iam.serviceAccountAdmin
```
### 5. create new key credentials for the service account, download the credentials as a json file and move this file as `gcp-sa-credentials.json` into the terrform root directory where this README file resides. This credentials file is used in terraform configuration to provision GCP infrastructure
```
gcloud iam service-accounts keys create gcp-sa-credentials.json --iam-account=terraform-sockets-sa@bebo-dev-sockets-terraform.iam.gserviceaccount.com
```
### 6. required environment variables to be pre-configured and present that are used by terraform configuration

* TF_VAR_GCP_POSTGRES_PASSWORD. A super secret password for the postgres database postgres user account

### 7. initialise terraform and apply the configuration to GCP
```
terraform init

terraform apply
```