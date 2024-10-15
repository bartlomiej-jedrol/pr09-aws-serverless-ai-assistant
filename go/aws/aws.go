package aws

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
)

var (
	ErrorFailedToLoadEnvFile       error = errors.New("failed to load .env file")
	ErrorBlankEnvVar               error = errors.New("blank environment variable")
	ErrorFailedToLoadDefaultConfig error = errors.New("failed to load default config")
)

// LoadDefaultConfig load default AWS config.
func LoadDefaultConfig() (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("ERROR: LoadDefaultConfig - %v, %v", ErrorFailedToLoadDefaultConfig, err)
		return aws.Config{}, ErrorFailedToLoadDefaultConfig
	}
	return cfg, nil
}

// GetEnvironmentVariable returns environment variable's value.
func GetEnvironmentVariable(envVarName string) (*string, error) {
	log.Printf("INFO: GetEnvironmentVariable - Entering GetEnvironmentVariable")

	err := godotenv.Load()
	if err != nil {
		log.Printf("ERROR: GetEnvironmentVariable - %v, %v", ErrorFailedToLoadEnvFile, err)
		return nil, ErrorFailedToLoadEnvFile
	}

	ev := os.Getenv(envVarName)
	if ev == "" {
		log.Printf("ERROR: GetEnvironmentVariable - %v, %s", ErrorBlankEnvVar, envVarName)
		return nil, ErrorBlankEnvVar
	}
	return &ev, nil
}
