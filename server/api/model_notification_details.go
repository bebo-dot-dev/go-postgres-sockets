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

type NotificationDetails struct {

	// describes the type of notification (0 = none, 1 = email, 2 = sms, 3 = slack)
	NotificationType int32 `json:"notificationType"`

	// arbitrary notification data
	NotificationText string `json:"notificationText"`
}
