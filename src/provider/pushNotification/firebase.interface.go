package pushNotification

import "context"

type IPushService interface {
	SendToToken(ctx context.Context, token string, payload NotificationPayload) (*PushSendResponse, error)
	SendToMultipleTokens(ctx context.Context, tokens []string, payload NotificationPayload) (*PushBatchResponse, error)
	SendToTopic(ctx context.Context, topic string, payload NotificationPayload) (*PushSendResponse, error)
	SendDataOnlyMessage(ctx context.Context, token string, data map[string]string) (*PushSendResponse, error)
	SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*TopicManagementResponse, error)
	UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*TopicManagementResponse, error)
	ValidateToken(ctx context.Context, token string) error
	SendToCondition(ctx context.Context, condition string, payload NotificationPayload) (*PushSendResponse, error)
}
