package slack

import (
	"encoding/json"
	"errors"
	"log"
)

var (
	ErrorEmptyJSONForUnmarshaling error = errors.New("empty JSON provided for unmarshaling")
	ErrorFailedToUnmarshalJSON    error = errors.New("failed to unmarshal JSON")
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

// UnmarshalJSON unmarshals provided JSON.
func UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		log.Printf("ERROR: UnmarshalJSON - %v", ErrorEmptyJSONForUnmarshaling)
		return ErrorEmptyJSONForUnmarshaling
	}

	slackMessages := []SlackMessage{}
	err := json.Unmarshal(data, &slackMessages)
	if err != nil {
		log.Printf("ERROR: UnmarshalJSON - %v, %v", ErrorFailedToUnmarshalJSON, err)
		return err
	}
	return nil
}
