package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

const (
	dbport = 5432
	dbname = "notifications"
)

// NotificationsApiService is a service that implements the logic for the NotificationsApiServicer
// This service should implement the business logic for every endpoint for the NotificationsAPIApi API.
// Include any external packages or services that will be required by this service.
type NotificationsApiService struct {
	apiKeys []string
}

// NewNotificationsApiService creates a default api service
func NewNotificationsApiService() NotificationsApiServicer {
	return &NotificationsApiService{}
}

func (s *NotificationsApiService) getDbConnection() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		dbport,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func (s *NotificationsApiService) isApiKeyAuthenticated(apiKey string) bool {
	for _, k := range s.apiKeys {
		if k == apiKey {
			return true
		}
	}
	return false
}

// Authenticate a user request
func (s *NotificationsApiService) Authenticate(ctx context.Context, request AuthenticationDetails) (ImplResponse, error) {
	err := request.Validate()
	if err != nil {
		return Response(http.StatusBadRequest, nil), err
	}

	if request.AuthKey != os.Getenv("GO_POSTGRES_SOCKETS_AUTH_KEY") {
		return Response(http.StatusUnauthorized, nil), errors.New("incorrect AuthKey supplied")
	}

	// store a new generated apiKey and return it to the caller for use in subsequent api calls
	newApiKey := uuid.NewString()
	s.apiKeys = append(s.apiKeys, newApiKey)
	return Response(http.StatusAccepted, newApiKey), nil
}

// AddNotification - adds a new notification
func (s *NotificationsApiService) AddNotification(ctx context.Context, request NotificationDetails) (ImplResponse, error) {
	err := request.Validate()
	if err != nil {
		return Response(http.StatusBadRequest, nil), err
	}

	if !s.isApiKeyAuthenticated(request.ApiKey) {
		return Response(http.StatusUnauthorized, nil), errors.New("incorrect ApiKey supplied")
	}

	db := s.getDbConnection()
	defer db.Close()

	var notificationId sql.NullInt32
	err = db.QueryRow("SELECT public.new_notification($1, $2);", request.NotificationType, request.NotificationText).Scan(&notificationId)

	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(http.StatusCreated, Id{Id: notificationId.Int32}), nil
}

// Ping - tests this api
func (s *NotificationsApiService) Ping(ctx context.Context) (ImplResponse, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}
	return Response(http.StatusCreated, PingResponse{Hostname: hostname}), nil
}
