package models

type LinkShortenerInputPayload struct {
	LongLink string `json:"longLink"`
}

type LinkShortenerOutputPayload struct {
	ShortLink string `json:"shortLink"`
}
