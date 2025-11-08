package model

type WatchStreamRequest struct {
	ChannelID string `json:"channel_id"`
	Name      string `json:"name"`
}

type WatchStreamResponse struct {
	URL string `json:"url"`
}
