package firebase

import (
	"context"
	"errors"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"

	"github.com/birukbelay/gocmn/src/logger"
)

type FirebaseServ struct {
	FirebaseApp *firebase.App
	FCMClient   *messaging.Client
}

// InitializeFirebase initializes Firebase app and FCM clien// TO Be called ON main
func InitFirebaseWithServAccPath(serviceAccountPath string) (*FirebaseServ, error) {
	ctx := context.Background()

	// Get the service account key path from environment variable
	if serviceAccountPath == "" {
		log.Println("Warning: FIREBASE_SERVICE_ACCOUNT_PATH not set, using default credentials")
		return nil, errors.New("FIREBASE_SERVICE_ACCOUNT_PATH not set, using default credentials")
	}

	opt := option.WithAuthCredentialsFile(option.ServiceAccount, serviceAccountPath)

	// Initialize Firebase app

	LocFirebaseApp, err := firebase.NewApp(ctx, nil, opt)

	if err != nil {
		logger.LogTrace("Error initializing Firebase app:", err)
		return nil, err
	}

	// Initialize FCM client
	FCMClient, err := LocFirebaseApp.Messaging(ctx)
	if err != nil {
		logger.LogTrace("Error initializing FCM client: From Init Func", err)
		return nil, err
	}

	log.Println("Firebase and FCM initialized successfully")
	return &FirebaseServ{FirebaseApp: LocFirebaseApp, FCMClient: FCMClient}, nil
}

