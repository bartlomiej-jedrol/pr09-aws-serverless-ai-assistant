package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/pkg/models"
	"github.com/bartlomiej-jedrol/pr09-aws-serverless-ai-assistant/go/pkg/utils"
	dub "github.com/dubinc/dub-go"
	"github.com/dubinc/dub-go/models/operations"
)

func ShortenLinkWithDub(token string, longLink string) string {
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

func HandleRequest(ctx context.Context, event json.RawMessage) (models.LinkShortenerOutputPayload, error) {
	log.Printf("input from router_lambda to link_shortener lambda: %s", string(event))
	token := utils.GetEnvVariable("DUB_API_KEY")
	if token == nil {
		log.Println("failed to get dub api token")
		return models.LinkShortenerOutputPayload{}, errors.New("failed to get dub api token")
	}

	inputPayload := models.LinkShortenerInputPayload{}
	err := json.Unmarshal(event, &inputPayload)
	if err != nil {
		log.Printf("failed to unmarshal JSON: %s", err)
		return models.LinkShortenerOutputPayload{}, err
	}

	log.Printf("long link: %s", inputPayload.LongLink)
	outputPayload := models.LinkShortenerOutputPayload{}
	// longLink := "https://gist.github.com/bartlomiej-jedrol/6c1010ae182de054641608b04eecacfe"
	outputPayload.ShortLink = ShortenLinkWithDub(*token, inputPayload.LongLink)
	log.Printf("short link: %s", outputPayload.ShortLink)

	return outputPayload, nil
}

func main() {
	lambda.Start(HandleRequest)
}
