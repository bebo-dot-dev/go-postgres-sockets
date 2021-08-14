# go-postgres-sockets

## This is an outline set of steps that result in this application being built and hosted in a Google Cloud Platform (GCP) k8s Google Kubernetes Engine (GKE) cluster. These steps assume a local install of the GCP SDK and/or use of the Google Cloud Shell for executing `gcloud` and `kubectl` commands.
---

### 1. create an artifact repository in the desired target location to act as the target for the containerised application
```
gcloud artifacts repositories create go-postgres-sockets-repo \
  --project=gke-test-322018 \
  --repository-format=docker \
  --location=europe-west2 \
  --description="Docker repository"
```
### 2. build the container with cloud build. This is similar to running docker build and docker push, but in this case the build happens on GCP and the container gets hosted in the artifact repository previously created
```
gcloud builds submit \
  --tag europe-west2-docker.pkg.dev/gke-test-322018/go-postgres-sockets-repo/go-postgres-sockets
```
### 3. create a GKE cluster with workload identity enabled. A GKE cluster is a managed set of Compute Engine virtual machines that operate as a single GKE cluster. Workload identity enabled is required for this application because that is a Cloud SQL Proxy requirement (postgres db connection, more details below)
```
gcloud container clusters create go-postgres-sockets-gke \
  --num-nodes 1 \
  --zone europe-west2 \
  --workload-pool=gke-test-322018.svc.id.goog
```
### 4. create a new node pool in the new GKE cluster with workload identity enabled. A node pool is a group of nodes within a cluster that all have the same configuration. A node is a worker machine that run a containerized application and other workloads.
```
gcloud container node-pools create socket-pool \
  --cluster=go-postgres-sockets-gke \
  --num-nodes 1 \
  --zone europe-west2 \
  --workload-metadata=GKE_METADATA
```
### 5. configure kubectl to communicate with the cluster
```
gcloud container clusters get-credentials go-postgres-sockets-gke \
  --zone europe-west2
```
### 6. create a namespace for a new k8s service account. Important: from this point on, GCP commands and yaml configuration files will assume that the operating namespace for this application is in this namespace so it is necessary to use and include this namespace in GCP commands and yaml configuration files from this point onwards
```
kubectl create namespace go-postgres-sockets-gke-ns
```
### 7. create a new named k8s service account
```
kubectl create serviceaccount --namespace go-postgres-sockets-gke-ns go-postgres-sockets-gke-ksa
```
### 8. create a new named Google service account
```
gcloud iam service-accounts create go-postgres-sockets-gke-sa
```
### 9. allow the k8s service account to impersonate the Google service account by creating an IAM policy binding between the two
```
gcloud iam service-accounts add-iam-policy-binding \
  --role roles/iam.workloadIdentityUser \
  --member "serviceAccount:gke-test-322018.svc.id.goog[go-postgres-sockets-gke-ns/go-postgres-sockets-gke-ksa]" \
  go-postgres-sockets-gke-sa@gke-test-322018.iam.gserviceaccount.com
```
### 10. add a Google service account annotation to the k8s service account, using the email address of the Google service account
```
kubectl annotate serviceaccount \
  --namespace go-postgres-sockets-gke-ns \
  go-postgres-sockets-gke-ksa \
  iam.gke.io/gcp-service-account=go-postgres-sockets-gke-sa@gke-test-322018.iam.gserviceaccount.com
```
### 11. add required roles to the service account, cloudsql.client, cloudsql.admin, iam.workloadIdentityUser
```
gcloud projects add-iam-policy-binding gke-test-322018 \
  --member serviceAccount:go-postgres-sockets-gke-sa@gke-test-322018.iam.gserviceaccount.com \
  --role roles/cloudsql.client

gcloud projects add-iam-policy-binding gke-test-322018 \
  --member serviceAccount:go-postgres-sockets-gke-sa@gke-test-322018.iam.gserviceaccount.com \
  --role roles/cloudsql.admin

gcloud projects add-iam-policy-binding gke-test-322018 \
  --member serviceAccount:go-postgres-sockets-gke-sa@gke-test-322018.iam.gserviceaccount.com \
  --role roles/iam.workloadIdentityUser
```
### 12. check the roles assigned to the service account
```
gcloud projects get-iam-policy gke-test-322018  \
--flatten="bindings[].members" \
--format='table(bindings.role)' \
--filter="bindings.members:go-postgres-sockets-gke-sa@gke-test-322018.iam.gserviceaccount.com"
```
### 13. setup required secrets (environment variables) used by the application that are subseqently introduced in [deployment.yaml](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/deployment.yaml)
```
kubectl create secret generic gke-secrets \
  --namespace go-postgres-sockets-gke-ns \
  --from-literal=db-host=127.0.0.1 \
  --from-literal=postgres-user=postgres \
  --from-literal=postgres-password=<super_secret> \
  --from-literal=go-postgres-sockets-auth-key=<super_secret>
```
### 14. Deployment. deploy the built application container to the GKE cluster using the [deployment.yaml](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/deployment.yaml) configuration file. deployment.yaml also installs cloud-sql-proxy in "sidecar" mode alongside the application to enable the application to connect to the postgres db in a secure manner
```
kubectl apply -f deployment.yaml
```
### 15. create a new public static ip address for ingress to the k8s cluster. This new ip address will act as the single exposed public ip address to the k8s cluster
```
gcloud compute addresses create sockets-ingress-ip --global

# check the newly created ip address:

gcloud compute addresses describe sockets-ingress-ip --global --format='value(address)'
```
### 16. external, use a registrar DNS service to point a domain DNS A record at the new public static ip address
### 17. start provisioning of a new Google managed SSL certificate described in the [certificate.yaml](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/certificate.yaml) configuration file for the domain. The certificate is used in the HTTPS load balancer [ingress.yaml](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/ingress.yaml) configuration
```
kubectl apply -f certificate.yaml
```
### 18. create a new ClusterIP service described in the [service.yaml](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/service.yaml) configuration file pointing to the application listening port 8080
```
kubectl apply -f service.yaml
```
### 19. create the ingress to the cluster using the previously created public ip address and the previously created SSL certificate as described in the [ingress.yaml](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/ingress.yaml) configuration file. 
```
kubectl apply -f ingress.yaml
```
### 20. check the provisioning status of the SSL certificate on the ingress point. It can take up to one hour for the SSL certificate provisioning process to complete and for the certificate to become active
```
kubectl describe managedcertificate sockets-certificate --namespace go-postgres-sockets-gke-ns
```
