# go-postgres-sockets

### This directory contains the terraform configuration equivalent of what is implemented with native gcloud and kubectl cli commands in the [k8s](https://github.com/bebo-dot-dev/go-postgres-sockets/tree/main/cmd/go-postgres-sockets/k8s) directory. Applying the terraform configuration in this directory will result in this application being built and hosted in a Google Cloud Platform (GCP) k8s Google Kubernetes Engine (GKE) cluster at https://sockets.bebo.dev

### The terraform configuration in this directory assumes local installs of `cloud_sql_proxy` and `psql` to enable a restore of the [postgreg db](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/postgres/notifications_db_backup.sql)
---
### 1. create a new GCP project and set it as the default project. The GCP project Id needs to be unique across all projects ever created by anyone in GCP
```
gcloud projects create bebo-dev-sockets-terraform --name="go postgres sockets (terraform)"

gcloud config set project bebo-dev-sockets-terraform
```
### 2. enable billing on the new project. The billing accountId can be found in the GCP billing console
```
gcloud beta billing projects link bebo-dev-sockets-terraform --billing-account=xxxxxx-xxxxxx-xxxxxx
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
gcloud services enable sqladmin.googleapis.com
```
### 4. create a new named Google service account and make the new account an editor
```
gcloud iam service-accounts create terraform-sockets-sa

gcloud projects add-iam-policy-binding bebo-dev-sockets \
  --member serviceAccount:terraform-sockets-sa@bebo-dev-sockets-terraform.iam.gserviceaccount.com \
  --role roles/editor
```
### 5. create new key credentials for the service account and download the creds as a json file
```
gcloud iam service-accounts keys create bebo-dev-sockets-terraform-credentials.json --iam-account=terraform-sockets-sa@bebo-dev-sockets-terraform.iam.gserviceaccount.com
```