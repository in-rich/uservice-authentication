package config

import (
	"context"
	_ "embed"
	"github.com/in-rich/lib-go/deploy"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

//go:embed firebase-dev.yml
var firebaseDevFile []byte

//go:embed firebase-staging.yml
var firebaseStagingFile []byte

//go:embed firebase-prod.yml
var firebaseProdFile []byte

type firebaseConfig struct {
	AuthDomain        string `yaml:"auth-domain"`
	ProjectID         string `yaml:"project-id"`
	StorageBucket     string `yaml:"storage-bucket"`
	MessagingSenderID string `yaml:"messaging-sender-id"`
	AppID             string `yaml:"app-id"`
	MeasurementID     string `yaml:"measurement-id"`
	APIKey            string `yaml:"api-key"`
}

var Firebase = deploy.LoadConfig[firebaseConfig](
	deploy.ProdConfig(firebaseProdFile),
	deploy.StagingConfig(firebaseStagingFile),
	deploy.DevConfig(firebaseDevFile),
)

var AuthClient *auth.Client

func init() {
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: Firebase.ProjectID})
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	authApp, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to create auth client: %v", err)
	}

	AuthClient = authApp
}
