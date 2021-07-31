/*
 * User API
 *
 * A notifications proof of concept API
 *
 * API version: 0.1.0
 * Contact: joe@bebo.dev
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

import (
	"context"
	"net/http"
)

// NotificationsAPIApiRouter defines the required methods for binding the api requests to a responses for the NotificationsAPIApi
// The NotificationsAPIApiRouter implementation should parse necessary information from the http request,
// pass the data to a NotificationsAPIApiServicer to perform the required actions, then write the service results to the http response.
type NotificationsAPIApiRouter interface {
	AddNotification(http.ResponseWriter, *http.Request)
	Ping(http.ResponseWriter, *http.Request)
}

// NotificationsAPIApiServicer defines the api actions for the NotificationsAPIApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type NotificationsAPIApiServicer interface {
	AddNotification(context.Context, NotificationDetails) (ImplResponse, error)
	Ping(context.Context) (ImplResponse, error)
}
