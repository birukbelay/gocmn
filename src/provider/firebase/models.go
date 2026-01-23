package firebase

type FirebaseSendResponse struct {
	Success   bool
	MessageID string
	Error     error
}
type FirebaseBatchResponse struct {
	SuccessCount int
	FailureCount int
	Responses    []*FirebaseSendResponse
}

// NotificationPayload represents the structure of a push notification
type NotificationPayload struct {
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	ImageURL string            `json:"image_url,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

// ErrorInfo is a topic management error.
type ErrorInfo struct {
	Index  int
	Reason string
}
type TopicManagementResponse struct {
	SuccessCount int
	FailureCount int
	Errors       []*ErrorInfo
}
