// Main implements an entry point of the Lambda function.
package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/joho/godotenv"
)

var secret string

// init is called when the Lambda function is initialized.
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env file")
	}

	secretName := os.Getenv("SECRET_NAME")
	if secretName == "" {
		log.Printf("SECRET_NAME environment variable is blank")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load default config: %v", err)
	}

	secretsManagerClient := secretsmanager.NewFromConfig(cfg)
	secretValue, err := secretsManagerClient.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		log.Fatalf("Failed to get secret value: %v", err)
	}
	secret = aws.ToString(secretValue.SecretString)
	log.Printf("Secret: %s", secret)
}

// HandleRequest handles the request and returns a response.
func HandleRequest(
	request events.APIGatewayCustomAuthorizerRequest) (*events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken
	token = strings.TrimPrefix(token, "Bearer ")

	if token != secret {
		return &events.APIGatewayCustomAuthorizerResponse{
			PrincipalID: "user",
			PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
				Version: "2012-10-17",
				Statement: []events.IAMPolicyStatement{
					{
						Action:   []string{"execute-api:Invoke"},
						Effect:   "Deny",
						Resource: []string{request.MethodArn},
					},
				},
			},
		}, nil
	}

	return &events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "user",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Allow",
					Resource: []string{request.MethodArn},
				},
			},
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
