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

	token, err := awsInt.GetEnvironmentVariable("DUB_API_KEY")
	if err != nil {
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
func HandleRequest(ctx context.Context, request json.RawMessage) (slack.SlackResponse, error) {
	log.Printf("INFO: HandleRequest - handling link_shortener lambda event: %s", string(request))

	payload := []slack.SlackMessage{}
	err := json.Unmarshal(request, &payload)
	if err != nil {
		log.Printf("ERROR: HandleRequest - %v, %v", slack.ErrorFailedToUnmarshalJSON, err)
		return slack.SlackResponse{}, err
	}
	log.Printf("payload: %v", payload)

	var longLink string
	if len(payload) > 0 && len(payload[0].Elements) > 0 && len(payload[0].Elements[0].Elements) > 0 {
		for _, elem := range payload[0].Elements[0].Elements {
			element, ok := elem.(map[string]interface{})
			log.Printf("INFO: HandleRequest - element: %v", element)
			if ok && element["type"] == "link" {
				longLink = element["url"].(string)
				break
			}
		}
	}

	log.Printf("INFO: HandleRequest - long link: %s", longLink)

	shortLink := shortenLinkWithDub(dubToken, longLink)
	log.Printf("INFO: HandleRequest - short link: %s", shortLink)

	r := slack.SlackResponse{Response: shortLink}
	return r, nil
}

func main() {
	lambda.Start(HandleRequest)
}
