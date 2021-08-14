# go-postgres-sockets

### This is an outline set of steps that result in this application being built and hosted in a Google Cloud Platform (GCP) k8s Google Kubernetes Engine (GKE) cluster at https://sockets.bebo.dev

### These steps assume a local install of the GCP SDK and/or use of the Google Cloud Shell for executing `gcloud` and `kubectl` commands and local installs of `cloud_sql_proxy` and `psql` to enable a restore of the [postgreg db](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/postgres/notifications_db_backup.sql)
---

### 1. create a new GCP project and set it as the default project. The GCP project Id needs to be unique across all projects ever created by anyone in GCP
```
gcloud projects create bebo-dev-sockets --name="go postgres sockets"

gcloud config set project bebo-dev-sockets
```
### 2. enable billing on the new project. The billing accountId can be found in the GCP billing console
```
gcloud beta billing projects link bebo-dev-sockets --billing-account=xxxxxx-xxxxxx-xxxxxx
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
### 4. create a new cloud sql postgres instance and apply a new password to the postgres user. Creation of the new instance will take some time
```
gcloud sql instances create postgres \
--database-version=POSTGRES_13 \
--tier=db-f1-micro \
--region=europe-west2


gcloud sql users set-password postgres \
--instance=postgres \
--password=super_secret
```
### 5. retrieve the connection name of the new cloud sql instance
```
gcloud sql instances describe postgres | grep connectionName

# connectionName: bebo-dev-sockets:europe-west2:postgres
```
### 6. start `cloud_sql_proxy` on local and then restore the [notifications postgres database](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/postgres/notifications_db_backup.sql) with `psql`
```
./cloud_sql_proxy -instances=bebo-dev-sockets:europe-west2:postgres=tcp:0.0.0.0:1234

psql "host=127.0.0.1 port=1234 sslmode=disable dbname=postgres user=postgres password=super_secret" < notifications_db_backup.sql
```

### 7. create an artifact repository in the desired target location to act as the target for the containerised application
```
gcloud artifacts repositories create go-postgres-sockets-repo \
  --project=bebo-dev-sockets \
  --repository-format=docker \
  --location=europe-west2 \
  --description="Docker repository"
```
### 8. build the application container with cloud build. This is similar to running docker build and docker push, but in this case the build happens on GCP and the container gets hosted in the artifact repository previously created
```
# in the /cmd/go-postgres-sockets directory where the application and it's Dockerfile reside

gcloud builds submit \
  --tag europe-west2-docker.pkg.dev/bebo-dev-sockets/go-postgres-sockets-repo/go-postgres-sockets
```
### 9. create a new GKE cluster with workload identity enabled. A GKE cluster is a managed set of Compute Engine virtual machines that operate as a single GKE cluster. Workload identity enabled is required for this application because that's needed for `cloud_sql_proxy` and connecting to the postgres db (further details below)
```
gcloud container clusters create go-postgres-sockets-gke \
  --num-nodes 1 \
  --zone europe-west2 \
  --workload-pool=bebo-dev-sockets.svc.id.goog
```

### 10. configure kubectl to communicate with the cluster
```
gcloud container clusters get-credentials go-postgres-sockets-gke \
  --zone europe-west2
```
### 11. create a namespace for a new k8s service account. Important: from this point on, GCP commands and yaml configuration files will assume that the operating namespace for this application is in this namespace so it is necessary to use and include this namespace in GCP commands and yaml configuration files from this point onwards
```
kubectl create namespace go-postgres-sockets-gke-ns
```
### 12. create a new named k8s service account for running the application in workload identity
```
kubectl create serviceaccount --namespace go-postgres-sockets-gke-ns go-postgres-sockets-gke-ksa
```
### 13. create a new named Google service account
```
gcloud iam service-accounts create go-postgres-sockets-gke-sa
```
### 14. allow the k8s service account to impersonate the Google service account by creating an IAM policy binding between the two
```
gcloud iam service-accounts add-iam-policy-binding \
  --role roles/iam.workloadIdentityUser \
  --member "serviceAccount:bebo-dev-sockets.svc.id.goog[go-postgres-sockets-gke-ns/go-postgres-sockets-gke-ksa]" \
  go-postgres-sockets-gke-sa@bebo-dev-sockets.iam.gserviceaccount.com
```
### 15. add a Google service account annotation to the k8s service account, using the email address of the Google service account
```
kubectl annotate serviceaccount \
  --namespace go-postgres-sockets-gke-ns \
  go-postgres-sockets-gke-ksa \
  iam.gke.io/gcp-service-account=go-postgres-sockets-gke-sa@bebo-dev-sockets.iam.gserviceaccount.com
```
### 16. add the required roles to the service account: `roles/cloudsql.client` and `roles/iam.workloadIdentityUser`. These are the role priviledges that the account needs to run the application
```
gcloud projects add-iam-policy-binding bebo-dev-sockets \
  --member serviceAccount:go-postgres-sockets-gke-sa@bebo-dev-sockets.iam.gserviceaccount.com \
  --role roles/cloudsql.client

gcloud projects add-iam-policy-binding bebo-dev-sockets \
  --member serviceAccount:go-postgres-sockets-gke-sa@bebo-dev-sockets.iam.gserviceaccount.com \
  --role roles/iam.workloadIdentityUser
```
### 17. check the roles assigned to the service account
```
gcloud projects get-iam-policy bebo-dev-sockets  \
--flatten="bindings[].members" \
--format='table(bindings.role)' \
--filter="bindings.members:go-postgres-sockets-gke-sa@bebo-dev-sockets.iam.gserviceaccount.com"
```
### 18. setup required secrets (environment variables) used by the application that are subseqently introduced in [deployment.yaml](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/deployment.yaml)
```
kubectl create secret generic gke-secrets \
  --namespace go-postgres-sockets-gke-ns \
  --from-literal=db-host=127.0.0.1 \
  --from-literal=postgres-user=postgres \
  --from-literal=postgres-password=<super_secret> \
  --from-literal=go-postgres-sockets-auth-key=<super_secret>
```
### 19. Deployment. deploy the built application container to the GKE cluster using the [deployment.yaml](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/deployment.yaml) configuration file. deployment.yaml also installs cloud-sql-proxy in "sidecar" mode alongside the application to enable the application to connect to the postgres db in a secure manner
```
kubectl apply -f deployment.yaml
```
### 20. create a new public static ip address for ingress to the k8s cluster. This new ip address will act as the single exposed public ip address to the k8s cluster
```
gcloud compute addresses create sockets-ingress-ip --global

# check the newly created ip address:

gcloud compute addresses describe sockets-ingress-ip --global --format='value(address)'
```
### 21. external, use a registrar DNS service to point a domain DNS A record at the new public static ip address

### 22. start provisioning of a new Google managed SSL certificate described in the [certificate.yaml](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/certificate.yaml) configuration file for the domain. The certificate is used in the HTTPS load balancer [ingress.yaml](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/ingress.yaml) configuration
```
kubectl apply -f certificate.yaml
```
### 23. create a new ClusterIP service described in the [service.yaml](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/service.yaml) configuration file pointing to the application listening port 8080
```
kubectl apply -f service.yaml
```
### 24. create the ingress to the cluster using the previously created public ip address and the previously created SSL certificate as described in the [ingress.yaml](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/ingress.yaml) configuration file. 
```
kubectl apply -f ingress.yaml
```

### 25 since this is a websocket application, the default 30 second timeout values on the ingress backends isn't good enough, it kills persistent websocket connections in an application that does not implement keep alive ping/pong messages. It's not possible to specify a timeout value at service/ingress creation time so this update has to be applied following service ingress creation
```
# list the created back end names

gcloud compute backend-services list --format="value(selfLink.scope())"

# then update the timeout field for every backend that needs a higher timeout

gcloud compute backend-services update $BACKEND1 --global --timeout=86400
gcloud compute backend-services update $BACKEND2 --global --timeout=86400
```

### 26. check the provisioning status of the SSL certificate on the ingress point. It can take up to an hour for the SSL certificate provisioning process to complete and for the certificate to become active
```
kubectl describe managedcertificate sockets-certificate --namespace go-postgres-sockets-gke-ns
```

### 27. done, the application is now hosted in the new GCP k8s cluster and exposed on https://sockets/bebo.dev