package service

import (
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"gitlab.yurtal.tech/company/blitz/business-card/back/internal/config"
	"gitlab.yurtal.tech/company/blitz/business-card/back/internal/model"
)

// Note: strings, core, and config are used in GetFeaturedChannels and buildChannelResponse methods

// RedisClientI interface for Redis operations
type RedisClientI interface {
	GenerateURLToken(url string) (string, error)
	GetURLByToken(token string) (string, error)
	DeleteToken(token string) error
}

type Stream struct {
	app         *pocketbase.PocketBase
	redisClient RedisClientI
}

func NewStream(app *pocketbase.PocketBase, redisClient RedisClientI) *Stream {
	return &Stream{
		app:         app,
		redisClient: redisClient,
	}
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

	actualURL := record.GetString("url")
	if actualURL == "" {
		return nil, fmt.Errorf("channel url is empty")
	}

	// Generate token for the URL
	token, err := s.redisClient.GenerateURLToken(actualURL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate URL token: %w", err)
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
		URL:      token,
		Quality:  qualityValue,
		Title:    title,
		Logo:     logo,
		Category: category,
		Country:  country,
		Language: language,
	}, nil
}

// buildChannelResponse is a helper method to build a WatchStreamResponse from a channel record
func (s *Stream) buildChannelResponse(record *core.Record) *model.WatchStreamResponse {
	channel := record.GetString("channel")
	actualURL := record.GetString("url")
	title := record.GetString("title")

	// Generate token for the URL
	token, err := s.redisClient.GenerateURLToken(actualURL)
	if err != nil {
		// If Redis fails, we need to handle this - don't return empty token
		// Log the error and return the actual URL as fallback
		// This ensures playback works even if Redis is temporarily unavailable
		token = actualURL
	}

	// Return token instead of actual URL
	url := token

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
	}
}

// GetFeaturedChannels retrieves featured channels from the database
func (s *Stream) GetFeaturedChannels() ([]*model.WatchStreamResponse, error) {
	// Query the featured table to get all featured records
	featuredRecords, err := s.app.FindRecordsByFilter(
		"featured",
		"",         // empty filter to get all records
		"-created", // sort by created date, newest first
		0,          // limit (0 = no limit)
		0,          // offset
	)
	if err != nil {
		return nil, err
	}

	var responses []*model.WatchStreamResponse

	for _, featuredRecord := range featuredRecords {
		// Get the channel ID from the relation field
		channelID := featuredRecord.GetString("channel")
		if channelID == "" {
			continue
		}

		// Fetch the actual channel record
		channelRecord, err := s.app.FindRecordById("channels", channelID)
		if err != nil || channelRecord == nil {
			continue
		}

		response := s.buildChannelResponse(channelRecord)
		responses = append(responses, response)
	}

	return responses, nil
}

// GetChannelByName retrieves a single channel by its name
func (s *Stream) GetChannelByName(channelName string) (*model.WatchStreamResponse, error) {
	record, err := s.app.FindFirstRecordByFilter("channels", fmt.Sprintf("channel = '%s'", channelName))
	if err != nil || record == nil {
		return nil, nil
	}

	response := s.buildChannelResponse(record)
	return response, nil
}

// GetChannelsByCategory retrieves 12 channels from a specific category with good quality
// Excludes featured channels. If category is "All", returns best quality channels from any category.
func (s *Stream) GetChannelsByCategory(categoryName string) ([]*model.WatchStreamResponse, error) {
	cfg := config.GetConfig()
	featuredIDs := strings.Split(cfg.FeaturedChannels, ",")

	// Build filter to exclude featured channels
	var excludeFilters []string
	for _, id := range featuredIDs {
		id = strings.TrimSpace(id)
		if id != "" {
			excludeFilters = append(excludeFilters, fmt.Sprintf("id != '%s'", id))
		}
	}

	var filter string

	// If category is "All" or "all", just get all channels
	if strings.ToLower(categoryName) == "all" || categoryName == "" {
		if len(excludeFilters) > 0 {
			filter = strings.Join(excludeFilters, " && ")
		} else {
			filter = "" // No filter needed
		}
	} else {
		// Get category ID by name
		categoryRecord, err := s.app.FindFirstRecordByFilter("categories", fmt.Sprintf("name_1 = '%s'", categoryName))
		if err != nil || categoryRecord == nil {
			return []*model.WatchStreamResponse{}, nil
		}

		categoryID := categoryRecord.Id

		// Build the filter
		filter = fmt.Sprintf("category = '%s'", categoryID)
		if len(excludeFilters) > 0 {
			filter = fmt.Sprintf("%s && (%s)", filter, strings.Join(excludeFilters, " && "))
		}
	}

	// Get channels with good quality (prioritize higher quality)
	records, err := s.app.FindRecordsByFilter(
		"channels",
		filter,
		"-quality", // Sort by quality descending
		12,
		0,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find channels by category: %w", err)
	}

	var responses []*model.WatchStreamResponse
	for _, record := range records {
		response := s.buildChannelResponse(record)
		responses = append(responses, response)
	}

	return responses, nil
}

// GetRecommendedChannels retrieves 4 similar channels based on the current watching channel
// Priority: 1) Same language + same category, 2) Same language (any category), 3) Same category (any language)
// Always prioritizes best quality and excludes the watching channel itself.
func (s *Stream) GetRecommendedChannels(req *model.RecommendStreamRequest) ([]*model.WatchStreamResponse, error) {
	var allResponses []*model.WatchStreamResponse

	// Get language and category IDs
	var languageID, categoryID string

	if req.LanguageName != "" {
		languageRecord, err := s.app.FindFirstRecordByFilter("languages", fmt.Sprintf("name = '%s'", req.LanguageName))
		if err == nil && languageRecord != nil {
			languageID = languageRecord.Id
		}
	}

	if req.CategoryName != "" {
		categoryRecord, err := s.app.FindFirstRecordByFilter("categories", fmt.Sprintf("name_1 = '%s'", req.CategoryName))
		if err == nil && categoryRecord != nil {
			categoryID = categoryRecord.Id
		}
	}

	// Strategy 1: Same language + same category
	if languageID != "" && categoryID != "" {
		filter := fmt.Sprintf("channel != '%s' && language = '%s' && category = '%s'", req.Channel, languageID, categoryID)
		records, err := s.app.FindRecordsByFilter("channels", filter, "-quality", 4, 0, nil)
		if err == nil {
			for _, record := range records {
				allResponses = append(allResponses, s.buildChannelResponse(record))
			}
		}
	}

	// If we have 4 channels, return them
	if len(allResponses) >= 4 {
		return allResponses[:4], nil
	}

	// Strategy 2: Same language (any category) - if we need more channels
	if languageID != "" && len(allResponses) < 4 {
		filter := fmt.Sprintf("channel != '%s' && language = '%s'", req.Channel, languageID)
		needed := 4 - len(allResponses)
		records, err := s.app.FindRecordsByFilter("channels", filter, "-quality", needed+10, 0, nil)
		if err == nil {
			// Add channels that we don't already have
			existingChannels := make(map[string]bool)
			for _, resp := range allResponses {
				existingChannels[resp.Channel] = true
			}

			for _, record := range records {
				channelName := record.GetString("channel")
				if !existingChannels[channelName] && len(allResponses) < 4 {
					allResponses = append(allResponses, s.buildChannelResponse(record))
					existingChannels[channelName] = true
				}
			}
		}
	}

	// If we still need more, Strategy 3: Same category (any language)
	if categoryID != "" && len(allResponses) < 4 {
		filter := fmt.Sprintf("channel != '%s' && category = '%s'", req.Channel, categoryID)
		needed := 4 - len(allResponses)
		records, err := s.app.FindRecordsByFilter("channels", filter, "-quality", needed+10, 0, nil)
		if err == nil {
			// Add channels that we don't already have
			existingChannels := make(map[string]bool)
			for _, resp := range allResponses {
				existingChannels[resp.Channel] = true
			}

			for _, record := range records {
				channelName := record.GetString("channel")
				if !existingChannels[channelName] && len(allResponses) < 4 {
					allResponses = append(allResponses, s.buildChannelResponse(record))
					existingChannels[channelName] = true
				}
			}
		}
	}

	// If we still don't have enough, get any high-quality channels
	if len(allResponses) < 4 {
		filter := fmt.Sprintf("channel != '%s'", req.Channel)
		needed := 4 - len(allResponses)
		records, err := s.app.FindRecordsByFilter("channels", filter, "-quality", needed+10, 0, nil)
		if err == nil {
			existingChannels := make(map[string]bool)
			for _, resp := range allResponses {
				existingChannels[resp.Channel] = true
			}

			for _, record := range records {
				channelName := record.GetString("channel")
				if !existingChannels[channelName] && len(allResponses) < 4 {
					allResponses = append(allResponses, s.buildChannelResponse(record))
					existingChannels[channelName] = true
				}
			}
		}
	}

	return allResponses, nil
}

// GetAllStreams retrieves all streams with filtering by category, country, language
// Returns paginated results (24 per page) sorted by quality
func (s *Stream) GetAllStreams(req *model.AllStreamsRequest) (*model.AllStreamsResponse, error) {
	const perPage = 24

	var filters []string

	// Filter by category if not "all"
	if req.Category != "" && strings.ToLower(req.Category) != "all" {
		categoryRecord, err := s.app.FindFirstRecordByFilter("categories", fmt.Sprintf("name_1 = '%s'", strings.ToLower(req.Category)))
		if err == nil && categoryRecord != nil {
			filters = append(filters, fmt.Sprintf("category = '%s'", categoryRecord.Id))
		}
	}

	// Filter by country if not "all"
	if req.Country != "" && strings.ToLower(req.Country) != "all" {
		countryRecord, err := s.app.FindFirstRecordByFilter("countries", fmt.Sprintf("name = '%s'", req.Country))
		if err == nil && countryRecord != nil {
			filters = append(filters, fmt.Sprintf("country = '%s'", countryRecord.Id))
		}
	}

	// Filter by language if not "all"
	if req.Language != "" && strings.ToLower(req.Language) != "all" {
		languageRecord, err := s.app.FindFirstRecordByFilter("languages", fmt.Sprintf("name = '%s'", req.Language))
		if err == nil && languageRecord != nil {
			filters = append(filters, fmt.Sprintf("language = '%s'", languageRecord.Id))
		}
	}

	// Build the filter
	filter := strings.Join(filters, " && ")

	// First, get total count
	totalRecords, err := s.app.FindRecordsByFilter("channels", filter, "", 0, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to count channels: %w", err)
	}
	total := len(totalRecords)

	// Calculate pagination
	totalPages := (total + perPage - 1) / perPage
	if req.Page > totalPages && totalPages > 0 {
		req.Page = totalPages
	}
	offset := (req.Page - 1) * perPage

	// Get paginated channels sorted by quality (highest to lowest)
	records, err := s.app.FindRecordsByFilter(
		"channels",
		filter,
		"-quality", // Sort by quality descending
		perPage,
		offset,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find channels: %w", err)
	}

	var channels []*model.WatchStreamResponse
	for _, record := range records {
		response := s.buildChannelResponse(record)
		channels = append(channels, response)
	}

	return &model.AllStreamsResponse{
		Channels:   channels,
		Total:      total,
		Page:       req.Page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// GetCategories retrieves all unique categories from database
func (s *Stream) GetCategories() ([]string, error) {
	records, err := s.app.FindRecordsByFilter("categories", "", "name_1", 0, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	// Use map to track unique categories
	uniqueMap := make(map[string]bool)
	var categories []string

	for _, record := range records {
		name := record.GetString("name_1")
		if name != "" && !uniqueMap[name] {
			uniqueMap[name] = true
			categories = append(categories, name)
		}
	}

	return categories, nil
}

// GetCountries retrieves all unique countries from database
func (s *Stream) GetCountries() ([]string, error) {
	records, err := s.app.FindRecordsByFilter("countries", "", "name", 0, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch countries: %w", err)
	}

	// Use map to track unique countries
	uniqueMap := make(map[string]bool)
	var countries []string

	for _, record := range records {
		name := record.GetString("name")
		if name != "" && !uniqueMap[name] {
			uniqueMap[name] = true
			countries = append(countries, name)
		}
	}

	return countries, nil
}

// GetLanguages retrieves all unique languages from database
func (s *Stream) GetLanguages() ([]string, error) {
	records, err := s.app.FindRecordsByFilter("languages", "", "name", 0, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch languages: %w", err)
	}

	// Use map to track unique languages
	uniqueMap := make(map[string]bool)
	var languages []string

	for _, record := range records {
		name := record.GetString("name")
		if name != "" && !uniqueMap[name] {
			uniqueMap[name] = true
			languages = append(languages, name)
		}
	}

	return languages, nil
}

// SearchStreams searches for channels by title (case-insensitive)
func (s *Stream) SearchStreams(req *model.SearchStreamRequest) (*model.SearchStreamResponse, error) {
	if req.Query == "" {
		return &model.SearchStreamResponse{
			Channels: []*model.WatchStreamResponse{},
			Total:    0,
		}, nil
	}

	// Search by title field with case-insensitive partial matching
	// Using ?~ for case-insensitive regex matching in PocketBase
	escapedQuery := strings.ReplaceAll(req.Query, "'", "\\'")
	filter := fmt.Sprintf("(title ?~ '%s' || id ?~ '%s')",
		escapedQuery,
		escapedQuery)

	// Parameters: collection, filter, sort, limit, offset, params
	records, err := s.app.FindRecordsByFilter("channels", filter, "-quality", 20, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search channels: %w", err)
	}

	var channels []*model.WatchStreamResponse
	for _, record := range records {
		response := &model.WatchStreamResponse{
			Channel: record.GetString("id"),
			Title:   record.GetString("title"),
			URL:     record.GetString("url"),
		}

		// Expand quality
		if qualityID := record.GetString("quality"); qualityID != "" {
			if qualityRecord, err := s.app.FindRecordById("qualities", qualityID); err == nil {
				response.Quality = qualityRecord.GetString("name")
			}
		}

		// Expand logo
		if logoID := record.GetString("logo"); logoID != "" {
			if logoRecord, err := s.app.FindRecordById("logos", logoID); err == nil {
				response.Logo = &model.Logo{
					URL: logoRecord.GetString("logo_url"),
				}
			}
		}

		// Expand category
		if categoryID := record.GetString("category"); categoryID != "" {
			if categoryRecord, err := s.app.FindRecordById("categories", categoryID); err == nil {
				response.Category = &model.Category{
					Name1: categoryRecord.GetString("name_1"),
					Name2: categoryRecord.GetString("name_2"),
				}
			}
		}

		// Expand country
		if countryID := record.GetString("country"); countryID != "" {
			if countryRecord, err := s.app.FindRecordById("countries", countryID); err == nil {
				response.Country = &model.Country{
					Name: countryRecord.GetString("name"),
				}
			}
		}

		// Expand language
		if languageID := record.GetString("language"); languageID != "" {
			if languageRecord, err := s.app.FindRecordById("languages", languageID); err == nil {
				response.Language = &model.Language{
					Name: languageRecord.GetString("name"),
				}
			}
		}

		channels = append(channels, response)
	}

	return &model.SearchStreamResponse{
		Channels: channels,
		Total:    len(channels),
	}, nil
}

// PlayStream resolves a token to the actual stream URL, or fetches URL by channel_id, or returns direct URL
func (s *Stream) PlayStream(req *model.PlayStreamRequest) (*model.PlayStreamResponse, error) {
	// If URL is provided directly, return it
	if req.URL != "" {
		return &model.PlayStreamResponse{
			URL: req.URL,
		}, nil
	}

	// If channel_id is provided, fetch the URL from database
	if req.ChannelID != "" {
		filter := fmt.Sprintf("id = '%s'", req.ChannelID)
		record, err := s.app.FindFirstRecordByFilter("channels", filter)
		if err != nil || record == nil {
			return nil, fmt.Errorf("channel not found")
		}

		actualURL := record.GetString("url")
		if actualURL == "" {
			return nil, fmt.Errorf("channel url is empty")
		}

		return &model.PlayStreamResponse{
			URL: actualURL,
		}, nil
	}

	// If token is provided, try to resolve it from Redis
	if req.Token != "" {
		// First, try to get the actual URL from Redis using the token
		actualURL, err := s.redisClient.GetURLByToken(req.Token)
		if err == nil {
			// Token found in Redis, return the URL
			return &model.PlayStreamResponse{
				URL: actualURL,
			}, nil
		}

		// If token is not in Redis, check if it's a direct URL (fallback for Redis failures)
		// Direct URLs typically start with http:// or https://
		if strings.HasPrefix(req.Token, "http://") || strings.HasPrefix(req.Token, "https://") {
			return &model.PlayStreamResponse{
				URL: req.Token,
			}, nil
		}

		// Token not found and not a direct URL
		return nil, fmt.Errorf("invalid or expired token")
	}

	return nil, fmt.Errorf("token, channel_id, or url is required")
}
