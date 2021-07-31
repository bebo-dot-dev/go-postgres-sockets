# go-postgres-sockets

A proof of concept Go application that demonstrates the use of postgres NOTIFY (plpgsql pg_notify function) and LISTEN (Go postgres package *pq.Listener) in conjunction with https://github.com/gorilla/websocket. 

In this POC, data changes applied within a postgres db table result in Go listening code receiving notification of the data change and the data change then broadcast to connected websocket clients.