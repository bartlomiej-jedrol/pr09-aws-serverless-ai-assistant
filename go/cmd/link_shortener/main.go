package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	awsInt "github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/aws"
	"github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/slack"
	dub "github.com/dubinc/dub-go"
	"github.com/dubinc/dub-go/models/operations"
)

var dubToken string

// init is called when the Lambda function is initialized.
func init() {
	log.Println("INFO: init - initializing link_shortener lambda")

	envVarName := "DUB_API_KEY"
	token, err := awsInt.GetEnvironmentVariable(envVarName)
	if err != nil {
		log.Printf("ERROR: init - %v, %v, %v", awsInt.ErrorFailedToGetEnvVar, envVarName, err)
	}
	if token == nil {
		log.Printf("ERROR: init - %v, %v", awsInt.ErrorFailedToGetEnvVar, envVarName)
		return
	}
	dubToken = *token
}

// shortenLinkWithDub shortens provided link
func shortenLinkWithDub(token string, longLink string) string {
	d := dub.New(dub.WithSecurity(token))

	req := &operations.CreateLinkRequestBody{
		URL: longLink,
	}

	ctx := context.Background()
	res, err := d.Links.Create(ctx, req)
	if err != nil {
		log.Printf("ERROR: shortenLinkWithDub - failed to shorten link with Dub: %v", err)
	}
	return res.ShortLink
}

// HandleRequest handler for link_shortner lambda.
func HandleRequest(ctx context.Context, request json.RawMessage) (slack.LinkShortenerOutputPayload, error) {
	log.Printf("INFO: HandleRequest - handling link_shortener lambda event: %s", string(request))

	inputPayload := slack.LinkShortenerInputPayload{}
	err := json.Unmarshal(request, &inputPayload)
	if err != nil {
		log.Printf("ERROR: HandleRequest - %v, %v", slack.ErrorFailedToUnmarshalJSON, err)
		return slack.LinkShortenerOutputPayload{}, err
	}

	log.Printf("INFO: HandleRequest - long link: %s", inputPayload.LongLink)
	outputPayload := slack.LinkShortenerOutputPayload{}
	outputPayload.ShortLink = shortenLinkWithDub(dubToken, inputPayload.LongLink)
	log.Printf("INFO: HandleRequest - long link: %s", outputPayload.ShortLink)
	return outputPayload, nil
}

func main() {
	lambda.Start(HandleRequest)
}
