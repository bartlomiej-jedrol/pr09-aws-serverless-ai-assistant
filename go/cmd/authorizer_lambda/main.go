// Main implements an entry point of the Lambda function.
package main

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	awsInt "github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/aws"
)

var secret string

// init is called when the Lambda function is initialized.
func init() {
	log.Println("INFO: init - initializing authorizer_lambda")

	ev, _ := awsInt.GetEnvironmentVariable("SECRET_NAME")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("ERROR: init - Failed to load default config: %v", err)
	}

	secretsManagerClient := secretsmanager.NewFromConfig(cfg)
	secretValue, err := secretsManagerClient.GetSecretValue(
		context.TODO(), &secretsmanager.GetSecretValueInput{
			SecretId: ev,
		})
	if err != nil {
		log.Printf("ERROR: init - failed to get secret value from AWS Secret Manager: %v", err)
	}
	secret = aws.ToString(secretValue.SecretString)
}

// HandleRequest handles the request and returns a response.
func HandleRequest(request events.APIGatewayCustomAuthorizerRequest) (
	*events.APIGatewayCustomAuthorizerResponse, error) {
	log.Printf("INFO: HandleRequest - handling authorizer_lambda event: %v", request)

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
