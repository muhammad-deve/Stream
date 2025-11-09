package model

type WatchStreamRequest struct {
	ChannelID string `json:"channel_id"`
	Name      string `json:"name"`
}

type Logo struct {
	URL    string  `json:"url"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Category struct {
	Name1 string `json:"name_1"`
	Name2 string `json:"name_2"`
	Name3 string `json:"name_3"`
}

type Country struct {
	Name string `json:"name"`
}

type Language struct {
	Name string `json:"name"`
}

type WatchStreamResponse struct {
	Channel  string    `json:"channel"`
	Title    string    `json:"title"`
	URL      string    `json:"url"`
	Quality  string    `json:"quality"`
	Logo     *Logo     `json:"logo"`
	Category *Category `json:"category"`
	Country  *Country  `json:"country"`
	Language *Language `json:"language"`
}
