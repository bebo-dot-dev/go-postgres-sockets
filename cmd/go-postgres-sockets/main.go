package main

import (
	"embed"
	"github.com/bebo-dot-dev/go-postgres-sockets/server/api"
	"github.com/bebo-dot-dev/go-postgres-sockets/server/database"
	"github.com/bebo-dot-dev/go-postgres-sockets/server/socket"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static
var embeddedFiles embed.FS

func main() {
	log.Printf("Server started")

	service := api.NewNotificationsApiService()
	controller := api.NewNotificationsApiController(service)

	router := api.NewRouter(controller)

	hub := registerSocketHandler(router)
	listener := database.NewPostgresDbListener(hub)
	listener.Listen()

	router.PathPrefix("/").Handler(http.FileServer(getEmbeddedFileSystem()))

	corsOrigins := handlers.AllowedOrigins([]string{"https://editor.swagger.io"})
	corsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(corsOrigins, corsHeaders, corsMethods)(router)))
}

func registerSocketHandler(router *mux.Router) *socket.Hub {
	hub := socket.NewSocketHub()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeSocket(hub, w, r)
	})

	return hub
}

func getEmbeddedFileSystem() http.FileSystem {
	f, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		panic(err)
	}
	return http.FS(f)
}