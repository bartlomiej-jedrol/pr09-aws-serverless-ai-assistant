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
	log.Println("INFO [init] Initializing link_shortener lambda")

	token, err := awsInt.GetEnvironmentVariable("DUB_API_KEY")
	if err != nil {
		log.Printf("ERROR [init] Failed to get DUB_API_KEY environment variable: %v", err)
	}

	if token == nil {
		log.Printf("ERROR [init] DUB_API_KEY environment variable is blank")
	}

	dubToken = *token

}

func shortenLinkWithDub(token string, longLink string) string {
	d := dub.New(dub.WithSecurity(token))

	req := &operations.CreateLinkRequestBody{
		URL: longLink,
	}

	ctx := context.Background()
	res, err := d.Links.Create(ctx, req)
	if err != nil {
		log.Fatalf("failed to create shorten link with Dub: %s", err)
	}
	return res.ShortLink
}

// HandleRequest handler for link_shortner lambda.
func HandleRequest(ctx context.Context, event json.RawMessage) (slack.LinkShortenerOutputPayload, error) {
	log.Printf("INFO [HandleRequest] Handling link_shortener lambda event: %s", string(event))

	inputPayload := slack.LinkShortenerInputPayload{}
	err := json.Unmarshal(event, &inputPayload)
	if err != nil {
		log.Printf("failed to unmarshal JSON: %s", err)
		return slack.LinkShortenerOutputPayload{}, err
	}

	log.Printf("long link: %s", inputPayload.LongLink)
	outputPayload := slack.LinkShortenerOutputPayload{}
	// longLink := "https://gist.github.com/bartlomiej-jedrol/6c1010ae182de054641608b04eecacfe"
	outputPayload.ShortLink = shortenLinkWithDub(dubToken, inputPayload.LongLink)
	log.Printf("short link: %s", outputPayload.ShortLink)

	return outputPayload, nil
}

func main() {
	lambda.Start(HandleRequest)
}
