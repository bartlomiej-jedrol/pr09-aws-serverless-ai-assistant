package slack

import (
	"encoding/json"
	"errors"
	"log"
)

type SlackResponse struct {
	Response string `json:"response"`
}

type Element struct {
	Type     string        `json:"type"`
	Elements []interface{} `json:"elements"`
}

type SlackMessage struct {
	Type     string    `json:"type"`
	BlockID  string    `json:"block_id"`
	Elements []Element `json:"elements"`
}

var (
	ErrorEmptyJSONForUnmarshaling error = errors.New("empty JSON provided for unmarshaling")
	ErrorFailedToUnmarshalJSON    error = errors.New("failed to unmarshal JSON")
	ErrorNoRequiredSlackElements  error = errors.New("slack payload does not have required elements")
	ErrorTypeAssertionFailed      error = errors.New("type assertion failed")
)

// UnmarshalSlackJSON unmarshals Slack message JSON.
func UnmarshalSlackJSON(request json.RawMessage) ([]SlackMessage, error) {
	payload := []SlackMessage{}
	if len(request) == 0 {
		log.Printf("ERROR: UnmarshalSlackJSON - %v", ErrorEmptyJSONForUnmarshaling)
		return nil, ErrorEmptyJSONForUnmarshaling
	}
	err := json.Unmarshal(request, &payload)
	if err != nil {
		log.Printf("ERROR: UnmarshalSlackJSON - %v, %v", ErrorFailedToUnmarshalJSON, err)
		return nil, ErrorFailedToUnmarshalJSON
	}
	log.Printf("payload: %v", payload)
	return payload, nil
}

// ValidElements validates if Slack message has required elements.
func ValidElements(payload []SlackMessage) bool {
	return len(payload) > 0 && len(payload[0].Elements) > 0 && len(payload[0].Elements[0].Elements) > 0
}

// ExtractElements extracts elements from Slack message.
func ExtractElements(payload []SlackMessage) (map[string]string, error) {
	elements := map[string]string{}
	if !ValidElements(payload) {
		log.Printf("ERROR: ExtractLink - %v", ErrorNoRequiredSlackElements)
		return elements, ErrorNoRequiredSlackElements
	}

	elements = map[string]string{}
	for _, elem := range payload[0].Elements[0].Elements {
		element, ok := elem.(map[string]interface{})
		log.Printf("INFO: ExtractElements - element: %v", element)
		if !ok {
			log.Printf("ERROR: ExtractElements - %v", ErrorTypeAssertionFailed)
			return elements, ErrorTypeAssertionFailed
		}
		if element["type"] == "text" {
			elements["text"], ok = element["text"].(string)
			if !ok {
				log.Printf("ERROR: ExtractElements - %v", ErrorTypeAssertionFailed)
				return elements, ErrorTypeAssertionFailed
			}
		}
		if element["type"] == "link" {
			elements["link"], ok = element["url"].(string)
			if !ok {
				log.Printf("ERROR: ExtractElements - %v", ErrorTypeAssertionFailed)
				return elements, ErrorTypeAssertionFailed
			}
		}
	}
	log.Printf("INFO: ExtractElements - elements: %v", elements)
	return elements, nil
}
