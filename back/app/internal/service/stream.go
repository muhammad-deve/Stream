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

	title := record.GetString("title")
	if title == "" {
		return nil, fmt.Errorf("channel title is empty")
	}

	channel := record.GetString("channel")
	
	// Get the quality value from the relation field
	qualityValue := ""
	qualityID := record.GetString("quality")
	if qualityID != "" {
		qualityRecord, err := s.app.FindRecordById("qualities", qualityID)
		if err == nil && qualityRecord != nil {
			qualityValue = qualityRecord.GetString("quality")
		}
	}

	// Get logo data
	var logo *model.Logo
	logoID := record.GetString("logo")
	if logoID != "" {
		logoRecord, err := s.app.FindRecordById("logos", logoID)
		if err == nil && logoRecord != nil {
			logo = &model.Logo{
				URL:    logoRecord.GetString("logo_url"),
				Width:  logoRecord.GetFloat("width"),
				Height: logoRecord.GetFloat("height"),
			}
		}
	}

	// Get category data
	var category *model.Category
	categoryID := record.GetString("category")
	if categoryID != "" {
		categoryRecord, err := s.app.FindRecordById("categories", categoryID)
		if err == nil && categoryRecord != nil {
			category = &model.Category{
				Name1: categoryRecord.GetString("name_1"),
				Name2: categoryRecord.GetString("name_2"),
				Name3: categoryRecord.GetString("name_3"),
			}
		}
	}

	// Get country data
	var country *model.Country
	countryID := record.GetString("country")
	if countryID != "" {
		countryRecord, err := s.app.FindRecordById("countries", countryID)
		if err == nil && countryRecord != nil {
			country = &model.Country{
				Name: countryRecord.GetString("name"),
			}
		}
	}

	// Get language data
	var language *model.Language
	languageID := record.GetString("language")
	if languageID != "" {
		languageRecord, err := s.app.FindRecordById("languages", languageID)
		if err == nil && languageRecord != nil {
			language = &model.Language{
				Name: languageRecord.GetString("name"),
			}
		}
	}

	return &model.WatchStreamResponse{
		Channel:  channel,
		URL:      url,
		Quality:  qualityValue,
		Title:    title,
		Logo:     logo,
		Category: category,
		Country:  country,
		Language: language,
	}, nil
}
