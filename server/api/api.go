package api

import (
	"context"
	"net/http"
)

// NotificationsApiRouter defines the required methods for binding the api requests to a responses for the NotificationsApi
// The NotificationsApiRouter implementation should parse necessary information from the http request,
// pass the data to a NotificationsApiServicer to perform the required actions, then write the service results to the http response.
type NotificationsApiRouter interface {
	Authenticate(http.ResponseWriter, *http.Request)
	AddNotification(http.ResponseWriter, *http.Request)
	Ping(http.ResponseWriter, *http.Request)
}

// NotificationsApiServicer defines the api actions for the NotificationsAPIApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type NotificationsApiServicer interface {
	Authenticate(context.Context, AuthenticationDetails) (ImplResponse, error)
	AddNotification(context.Context, NotificationDetails) (ImplResponse, error)
	Ping(context.Context) (ImplResponse, error)
}
