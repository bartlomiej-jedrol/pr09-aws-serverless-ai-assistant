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

var dubAPIKey string

// init initialises execution environment for AWS Lambda.
func init() {
	log.Println("INFO: init - initializing link_shortener lambda")
	apiKey, err := awsInt.GetEnvironmentVariable("DUB_API_KEY")
	if err != nil {
		return
	}
	dubAPIKey = *apiKey
}

// shortenLinkWithDub shortens provided link
func shortenLinkWithDub(apiKey string, longLink string) string {
	d := dub.New(dub.WithSecurity(apiKey))

	req := &operations.CreateLinkRequestBody{
		URL: longLink,
	}

	ctx := context.Background()
	resp, err := d.Links.Create(ctx, req)
	if err != nil {
		log.Printf("ERROR: shortenLinkWithDub - failed to shorten link with Dub: %v", err)
	}
	return resp.ShortLink
}

// HandleRequest handler for link_shortner lambda.
func HandleRequest(ctx context.Context, request json.RawMessage) (resp slack.SlackResponse, err error) {
	log.Printf("INFO: HandleRequest - handling link_shortener lambda event: %s", string(request))
	payload, err := slack.UnmarshalSlackJSON(request)
	if err != nil {
		return resp, err
	}
	elements, err := slack.ExtractElements(payload)
	if err != nil {
		return resp, err
	}
	link := elements["link"]

	shortLink := shortenLinkWithDub(dubAPIKey, link)
	log.Printf("INFO: HandleRequest - short link: %s", shortLink)
	resp = slack.SlackResponse{Response: shortLink}
	return resp, nil
}

func main() {
	lambda.Start(HandleRequest)
}
