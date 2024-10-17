// Package openai provides tools to work with openai models.
package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIPayload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

var (
	chatCompletionsEndpoint = "https://api.openai.com/v1/chat/completions"
)

// BuildPrompt returns..
func buildMessages(systemContent string, userContent string) []Message {
	messages := []Message{}
	messages = append(messages, Message{
		Role:    "system",
		Content: systemContent})
	messages = append(messages, Message{
		Role:    "user",
		Content: userContent})
	return messages
}

// BuildBody returns..
func buildBody(model string, systemContent string, userContent string) ([]byte, error) {
	messages := buildMessages(systemContent, userContent)
	payload := OpenAIPayload{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ERROR: buildBody - failed to marshal JSON, %v", err)
		return nil, err
	}
	return jsonData, nil
}

// sendRequest returns..
func sendRequest(enpointURL string, APIKey string, requestBody []byte) ([]byte, error) {
	requestDataReader := bytes.NewReader(requestBody)

	URL, err := url.Parse(enpointURL)
	if err != nil {
		log.Printf("ERROR: CreateChatCompletions - failed to parse URL, %v", err)
		return nil, err
	}

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", fmt.Sprintf("Bearer %s", APIKey))

	// log.Printf("header: %v", header.Get("Authorization"))
	// log.Printf("header: %v", header.Get("Content-Type"))

	request := http.Request{
		Method: http.MethodPost,
		URL:    URL,
		Header: header,
		Body:   io.NopCloser(requestDataReader),
	}

	log.Printf("Request Method: %s", request.Method)
	log.Printf("Request URL: %s", request.URL.String())
	log.Printf("Request Headers: %v", request.Header)
	log.Printf("Request Body: %s", string(requestBody))

	client := http.Client{}
	response, err := client.Do(&request)
	if err != nil {
		log.Printf("ERROR: sendRequest - failed to send request, %v", err)
		return nil, fmt.Errorf("ERROR: sendRequest - failed to send request, %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("ERROR: sendRequest - received non-OK HTTP status: %s", response.Status)
		return nil, fmt.Errorf("ERROR: sendRequest - received non-OK HTTP status: %s", response.Status)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("ERROR: sendRequest - failed to read response body, %v", err)
		return nil, err
	}
	return responseBody, nil

}

// CreateChatCompletions returns one or more predicted completions.
// It takes prompt and model, commmunicates with OpenAI models, and returns one or more predicted completions.
func CreateChatCompletions(openaiAPIKey string, model string, userContent string) error {
	systemContent := `Recognise a skill.
		Respond with JSON object. skill (link_shortener|other)
		`
	jsonRequestData, err := buildBody(model, systemContent, userContent)
	if err != nil {
		return err
	}

	jsonResponseData, err := sendRequest(chatCompletionsEndpoint, openaiAPIKey, jsonRequestData)
	if err != nil {
		return err
	}
	log.Printf("openai response: %s", string(jsonResponseData))
	return nil
}
