package api

type NotificationSocketMessage struct {

	// describes the database table changed
	Table string `json:"table"`

	// describes the database operation performed
	Operation string `json:"operation"`

	Data NotificationSocketMessageData `json:"data"`
}
