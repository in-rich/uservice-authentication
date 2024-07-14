package config

import (
	"context"
	_ "embed"
	"github.com/goccy/go-yaml"
	"os"

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

var Firebase *firebaseConfig

var AuthClient *auth.Client

func init() {
	cfg := new(firebaseConfig)

	switch os.Getenv("ENV") {
	case "prod":
		if err := yaml.Unmarshal(firebaseProdFile, cfg); err != nil {
			log.Fatalf("error unmarshalling firebase config: %v\n", err)
		}
	case "staging":
		if err := yaml.Unmarshal(firebaseStagingFile, cfg); err != nil {
			log.Fatalf("error unmarshalling firebase config: %v\n", err)
		}
	default:
		if err := yaml.Unmarshal(firebaseDevFile, cfg); err != nil {
			log.Fatalf("error unmarshalling firebase config: %v\n", err)
		}
	}

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: cfg.ProjectID})
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	authApp, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to create auth client: %v", err)
	}

	Firebase = cfg
	AuthClient = authApp
}
