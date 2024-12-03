package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	awsInt "github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/aws"
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
	dubAPIKey = apiKey
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
func HandleRequest(ctx context.Context, text json.RawMessage) (string, error) {
	log.Printf("INFO: HandleRequest - handling link_shortener lambda event: %s", string(text))

	t := string(text)
	log.Printf("text: %v", t)

	// shortLink := shortenLinkWithDub(dubAPIKey, t)
	// log.Printf("INFO: HandleRequest - short link: %s", shortLink)
	// resp = slack.SlackResponse{Response: shortLink}
	// return resp, nil
	return t, nil
}

func main() {
	lambda.Start(HandleRequest)
}
