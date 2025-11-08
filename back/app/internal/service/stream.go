package service

import (
	"fmt"

	"github.com/pocketbase/pocketbase"
	"gitlab.yurtal.tech/company/blitz/business-card/back/internal/model"
)

type Stream struct {
	app *pocketbase.PocketBase
}

func NewStream(app *pocketbase.PocketBase) *Stream {
	return &Stream{app: app}
}

func (s *Stream) WatchStream(req *model.WatchStreamRequest) (*model.WatchStreamResponse, error) {
	filter := fmt.Sprintf("id = '%s'", req.ChannelID)
	record, err := s.app.FindFirstRecordByFilter("channels", filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find channel: %w", err)
	}
	if record == nil {
		return nil, fmt.Errorf("channel not found")
	}

	url := record.GetString("url")
	if url == "" {
		return nil, fmt.Errorf("channel url is empty")
	}

	return &model.WatchStreamResponse{
		URL: url,
	}, nil
}
