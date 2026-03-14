package firebase

import (
	"context"
	"fmt"
	"os"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// Client is the initialized Firebase Auth client.
// It is safe for concurrent use and initialized exactly once.
var (
	authClient *auth.Client
	initErr    error
	once       sync.Once
)

// Init initializes the Firebase Admin SDK using the service account
// JSON file path specified in the FIREBASE_CREDENTIALS_PATH environment
// variable. It must be called once at application startup.
func Init(ctx context.Context) (*auth.Client, error) {
	once.Do(func() {
		credPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
		if credPath == "" {
			initErr = fmt.Errorf(
				"firebase.Init: FIREBASE_CREDENTIALS_PATH env var is not set",
			)
			return
		}
		if _, err := os.Stat(credPath); err != nil {
			initErr = fmt.Errorf(
				"firebase.Init: credentials file not found at %q: %w",
				credPath, err,
			)
			return
		}
		opt := option.WithCredentialsFile(credPath)
		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil {
			initErr = fmt.Errorf("firebase.Init: NewApp: %w", err)
			return
		}
		authClient, err = app.Auth(ctx)
		if err != nil {
			initErr = fmt.Errorf("firebase.Init: app.Auth: %w", err)
			return
		}
	})

	return authClient, initErr
}

// VerifyIDToken verifies the given Firebase ID token and returns
// the decoded token claims. Returns a wrapped error on failure.
func VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if authClient == nil {
		return nil, fmt.Errorf(
			"firebase.VerifyIDToken: client not initialized - call Init first",
		)
	}
	token, err := authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("firebase.VerifyIDToken: %w", err)
	}
	return token, nil
}
