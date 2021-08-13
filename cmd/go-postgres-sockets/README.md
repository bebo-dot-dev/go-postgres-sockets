# go-postgres-sockets

[main.go](https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/main.go) is a Go web application that uses the server [api, database and socket](https://github.com/bebo-dot-dev/go-postgres-sockets/tree/main/server) packages implemented in this repo. 

At the time of writing this application is hosted at https://sockets.bebo.dev in a Google Cloud Platform k8s Google Kubernetes Engine cluster.

Details of how this application is built and how the cluster and load balancer is provisioned are in https://github.com/bebo-dot-dev/go-postgres-sockets/blob/main/cmd/go-postgres-sockets/k8s/README.md