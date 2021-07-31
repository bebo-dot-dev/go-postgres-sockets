/*
 * Notifications API
 *
 * A notifications proof of concept API
 *
 * API version: 0.1.0
 * Contact: joe@bebo.dev
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

import "github.com/go-playground/validator"

type NotificationDetails struct {

	// describes the type of notification (0 = none, 1 = email, 2 = sms, 3 = slack)
	NotificationType int32 `json:"notificationType" validate:"gt=0,lte=3"`

	// arbitrary notification data
	NotificationText string `json:"notificationText" validate:"required"`
}

func (n *NotificationDetails) Validate() error {
	validate := validator.New()
	return validate.Struct(n)
}

