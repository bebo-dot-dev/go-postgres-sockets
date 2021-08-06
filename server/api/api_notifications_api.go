package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

// A NotificationsApiController binds http requests to an api service and writes the service results to the http response
type NotificationsApiController struct {
	service NotificationsApiServicer
}

// NewNotificationsApiController creates a default api controller
func NewNotificationsApiController(s NotificationsApiServicer) Router {
	return &NotificationsApiController{service: s}
}

// Routes returns all of the api route for the NotificationsApiController
func (c *NotificationsApiController) Routes() Routes {
	return Routes{
		{
			"Authenticate",
			strings.ToUpper("Post"),
			"/authenticate",
			c.Authenticate,
		},
		{
			"AddNotification",
			strings.ToUpper("Put"),
			"/addNotification",
			c.AddNotification,
		},
		{
			"Ping",
			strings.ToUpper("Get"),
			"/ping",
			c.Ping,
		},
	}
}

// Authenticate a user request
func (c *NotificationsApiController) Authenticate(w http.ResponseWriter, r *http.Request) {
	authenticationDetails := &AuthenticationDetails{}
	if err := json.NewDecoder(r.Body).Decode(&authenticationDetails); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := c.service.Authenticate(r.Context(), *authenticationDetails)
	// If an error occurred, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// AddNotification - adds a new notification
func (c *NotificationsApiController) AddNotification(w http.ResponseWriter, r *http.Request) {
	notificationDetails := &NotificationDetails{}
	if err := json.NewDecoder(r.Body).Decode(&notificationDetails); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := c.service.AddNotification(r.Context(), *notificationDetails)
	// If an error occurred, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// Ping - tests this api
func (c *NotificationsApiController) Ping(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.Ping(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
