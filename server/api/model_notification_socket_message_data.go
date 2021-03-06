package api

// NotificationSocketMessageData - notification data
type NotificationSocketMessageData struct {

	// the id of the row changed
	Id int32 `json:"id"`

	// describes the type of notification (0 = none, 1 = email, 2 = sms, 3 = slack)
	NotificationTypeId int32 `json:"notificationTypeId"`

	// the creation timestamp of the row
	CreatedTimestamp string `json:"createdTimestamp"`

	// notification text data
	NotificationText string `json:"notificationText"`

	// textual descriptor for notificationTypeId
	NotificationType string `json:"notificationType"`
}
