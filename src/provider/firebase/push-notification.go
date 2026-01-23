package firebase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"firebase.google.com/go/v4/messaging"

	"github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/util"
)

// FirebasePushService handles Firebase Cloud Messaging operations
type FirebasePushService struct {
	client *messaging.Client //the messaging client

}

// NewPushNotificationService creates a new push notification service instance
func (f *FirebaseServ) NewPushService() (*FirebasePushService, error) {

	if f.FirebaseApp == nil {
		return nil, errors.New("FCM client not initialized")
	}
	if f.FCMClient != nil {
		logger.LogTrace("fcm already initialized", "returning")
		return &FirebasePushService{
			client: f.FCMClient,
		}, nil
	}
	// Initialize FCM client
	NewFCMClient, err := f.FirebaseApp.Messaging(context.Background())
	if err != nil {
		logger.LogTrace("Error initializing FCM client :From service: %v", err)
		return nil, err
	}
	return &FirebasePushService{
		client: NewFCMClient,
	}, nil
}

// SendToToken sends a push notification to a specific device token
func (pns *FirebasePushService) SendToToken(ctx context.Context, token string, payload NotificationPayload) (*FirebaseSendResponse, error) {
	if pns.client == nil {
		return nil, fmt.Errorf("FCM client not initialized")
	}

	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title:    payload.Title,
			Body:     payload.Body,
			ImageURL: payload.ImageURL,
		},
		Data: payload.Data,
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Priority: messaging.PriorityHigh,
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: payload.Title,
						Body:  payload.Body,
					},
					Badge: intPtr(1),
				},
			},
		},
	}

	response, err := pns.client.Send(ctx, message)
	if err != nil {
		log.Printf("Error sending push notification to token %s: %v", token, err)
		return nil, err
	}

	log.Printf("Successfully sent push notification to token %s. Message ID: %s", token, response)
	return &FirebaseSendResponse{MessageID: response}, nil
}

// SendToMultipleTokens sends a push notification to multiple device tokens, must be upto 500
func (pns *FirebasePushService) SendToMultipleTokens(ctx context.Context, tokens []string, payload NotificationPayload) (*FirebaseBatchResponse, error) {
	if pns.client == nil {
		return nil, fmt.Errorf("FCM client not initialized")
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens provided")
	}

	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title:    payload.Title,
			Body:     payload.Body,
			ImageURL: payload.ImageURL,
		},
		Data: payload.Data,
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Priority: messaging.PriorityHigh,
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: payload.Title,
						Body:  payload.Body,
					},
					Badge: intPtr(1),
				},
			},
		},
	}

	response, err := pns.client.SendEachForMulticast(ctx, message)
	if err != nil {
		log.Printf("Error sending push notification to multiple tokens: %v", err)
		return nil, err
	}

	log.Printf("Successfully sent push notification to %d tokens. Success count: %d, Failure count: %d",
		len(tokens), response.SuccessCount, response.FailureCount)

	// Log failed tokens for debugging
	if response.FailureCount > 0 {
		for i, resp := range response.Responses {
			if !resp.Success {
				log.Printf("Failed to send to token %s: %v", tokens[i], resp.Error)
			}
		}
	}

	return &FirebaseBatchResponse{SuccessCount: response.SuccessCount, FailureCount: response.FailureCount,
		Responses: []*FirebaseSendResponse{{Success: true}},
	}, nil
}

// SendToTopic sends a push notification to a topic
func (pns *FirebasePushService) SendToTopic(ctx context.Context, topic string, payload NotificationPayload) (*FirebaseSendResponse, error) {
	if pns.client == nil {
		return nil, fmt.Errorf("FCM client not initialized")
	}

	message := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title:    payload.Title,
			Body:     payload.Body,
			ImageURL: payload.ImageURL,
		},
		Data: payload.Data,
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Priority: messaging.PriorityHigh,
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: payload.Title,
						Body:  payload.Body,
					},
					Badge: intPtr(1),
				},
			},
		},
	}

	response, err := pns.client.Send(ctx, message)
	if err != nil {
		log.Printf("Error sending push notification to topic %s: %v", topic, err)
		return nil, err
	}

	log.Printf("Successfully sent push notification to topic %s. Message ID: %s", topic, response)
	return &FirebaseSendResponse{MessageID: response}, nil
}

// SendDataOnlyMessage sends a data-only message (no notification UI)
func (pns *FirebasePushService) SendDataOnlyMessage(ctx context.Context, token string, data map[string]string) (*FirebaseSendResponse, error) {
	if pns.client == nil {
		return nil, fmt.Errorf("FCM client not initialized")
	}

	message := &messaging.Message{
		Token: token,
		Data:  data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
	}

	response, err := pns.client.Send(ctx, message)
	if err != nil {
		log.Printf("Error sending data-only message to token %s: %v", token, err)
		return nil, err
	}

	log.Printf("Successfully sent data-only message to token %s. Message ID: %s", token, response)
	return &FirebaseSendResponse{MessageID: response}, nil
}

// SubscribeToTopic subscribes device tokens to a topic
func (pns *FirebasePushService) SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*TopicManagementResponse, error) {
	if pns.client == nil {
		return nil, fmt.Errorf("FCM client not initialized")
	}

	response, err := pns.client.SubscribeToTopic(ctx, tokens, topic)
	if err != nil {
		log.Printf("Error subscribing tokens to topic %s: %v", topic, err)
		return nil, err
	}

	log.Printf("Successfully subscribed %d tokens to topic %s. Success count: %d, Failure count: %d",
		len(tokens), topic, response.SuccessCount, response.FailureCount)
	mapped, err := util.MarshalToStruct[TopicManagementResponse](response)
	if err != nil {
		log.Printf("Error subscribing tokens to topic %s: %v", topic, err)
		return nil, err
	}
	return mapped, nil
}

// UnsubscribeFromTopic unsubscribes device tokens from a topic
func (pns *FirebasePushService) UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*TopicManagementResponse, error) {
	if pns.client == nil {
		return nil, fmt.Errorf("FCM client not initialized")
	}

	response, err := pns.client.UnsubscribeFromTopic(ctx, tokens, topic)
	if err != nil {
		log.Printf("Error unsubscribing tokens from topic %s: %v", topic, err)
		return nil, err
	}

	log.Printf("Successfully unsubscribed %d tokens from topic %s. Success count: %d, Failure count: %d",
		len(tokens), topic, response.SuccessCount, response.FailureCount)

	mapped, err := util.MarshalToStruct[TopicManagementResponse](response)
	if err != nil {
		log.Printf("Error subscribing tokens to topic %s: %v", topic, err)
		return nil, err
	}
	return mapped, nil
}

// ValidateToken validates if a registration token is valid
func (pns *FirebasePushService) ValidateToken(ctx context.Context, token string) error {
	if pns.client == nil {
		return fmt.Errorf("FCM client not initialized")
	}

	// Send a test message with dry run to validate token
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: "Test",
			Body:  "Test message for token validation",
		},
	}

	_, err := pns.client.SendDryRun(ctx, message)
	if err != nil {
		log.Printf("Token validation failed for %s: %v", token, err)
		return err
	}

	log.Printf("Token %s is valid", token)
	return nil
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}

// SendToCondition sends a push notification based on topic conditions
func (pns *FirebasePushService) SendToCondition(ctx context.Context, condition string, payload NotificationPayload) (*FirebaseSendResponse, error) {
	if pns.client == nil {
		return nil, fmt.Errorf("FCM client not initialized")
	}

	// Validate condition before sending
	if err := validateCondition(condition); err != nil {
		return nil, fmt.Errorf("invalid condition: %v", err)
	}

	message := &messaging.Message{
		Condition: condition,
		Notification: &messaging.Notification{
			Title:    payload.Title,
			Body:     payload.Body,
			ImageURL: payload.ImageURL,
		},
		Data: payload.Data,
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Priority: messaging.PriorityHigh,
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: payload.Title,
						Body:  payload.Body,
					},
					Badge: intPtr(1),
				},
			},
		},
	}

	response, err := pns.client.Send(ctx, message)
	if err != nil {
		log.Printf("Error sending condition-based notification: %v", err)
		return nil, err
	}

	log.Printf("Successfully sent condition-based notification. Message ID: %s", response)
	return &FirebaseSendResponse{MessageID: response}, nil
}

// validateCondition validates FCM condition syntax
func validateCondition(condition string) error {
	// Count topics (max 5 allowed)
	topicCount := strings.Count(condition, "in topics")
	if topicCount > 5 {
		return fmt.Errorf("condition contains more than 5 topics (found %d)", topicCount)
	}

	if topicCount == 0 {
		return fmt.Errorf("condition must contain at least one topic")
	}

	// Check for balanced parentheses
	openCount := strings.Count(condition, "(")
	closeCount := strings.Count(condition, ")")
	if openCount != closeCount {
		return fmt.Errorf("unbalanced parentheses in condition")
	}

	// Basic syntax validation
	if !strings.Contains(condition, "in topics") {
		return fmt.Errorf("condition must contain 'in topics'")
	}

	return nil
}
