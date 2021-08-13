# go-postgres-sockets

https://sockets.bebo.dev (Google Cloud Platform [k8s hosted](https://github.com/bebo-dot-dev/go-postgres-sockets/tree/main/cmd/go-postgres-sockets/k8s))

A proof of concept Go application that demonstrates the use of postgres NOTIFY (plpgsql pg_notify function) and LISTEN (Go postgres package *pq.Listener) in conjunction with https://github.com/gorilla/websocket. 

In this POC, data changes applied within a postgres db table result in Go listening code receiving notification of the data change and the data change then broadcast to connected websocket clients.

## run / build
```
go run cmd/go-postgres-sockets main.go
```
```
go build cmd/go-postgres-sockets main.go
```

## video
![files](https://github.com/bebo-dot-dev/go-postgres-sockets/releases/download/0.0.1/go-postgres-sockets.gif)
