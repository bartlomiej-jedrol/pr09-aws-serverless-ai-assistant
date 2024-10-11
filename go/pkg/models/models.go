package models

type LinkShortenerInputPayload struct {
	LongLink string `json:"longLink"`
}

type LinkShortenerOutputPayload struct {
	ShortLink string `json:"shortLink"`
}

type TextElement struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type LinkElement struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type RichTextSection struct {
	Type     string `json:"type"`
	Elements []interface{}
}

type SlackMessage struct {
	Type     string `json:"type"`
	BlockID  string `json:"block_id"`
	Elements RichTextSection
}
