package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"gitlab.yurtal.tech/company/blitz/business-card/back/internal/model"
)

func (h *Handler) WatchStreamHandler(e *core.RequestEvent) error {
	var req model.WatchStreamRequest
	if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.ChannelID == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "channel_id is required",
		})
	}

	resp, err := h.service.Stream().WatchStream(&req)
	if err != nil {
		h.logger.Error("failed to watch stream", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, resp)
}

func (h *Handler) FeaturedStreamHandler(e *core.RequestEvent) error {
	resp, err := h.service.Stream().GetFeaturedChannels()
	if err != nil {
		h.logger.Error("failed to get featured channels", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, resp)
}

func (h *Handler) GetChannelHandler(e *core.RequestEvent) error {
	channelName := e.Request.PathValue("name")
	if channelName == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "channel name is required",
		})
	}

	resp, err := h.service.Stream().GetChannelByName(channelName)
	if err != nil {
		h.logger.Error("failed to get channel", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if resp == nil {
		return e.JSON(http.StatusNotFound, map[string]string{
			"error": "channel not found",
		})
	}

	return e.JSON(http.StatusOK, resp)
}

func (h *Handler) CategoryStreamHandler(e *core.RequestEvent) error {
	var req model.CategoryStreamRequest
	if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.CategoryName == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "category_name is required",
		})
	}

	resp, err := h.service.Stream().GetChannelsByCategory(req.CategoryName)
	if err != nil {
		h.logger.Error("failed to get channels by category", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, resp)
}

func (h *Handler) RecommendStreamHandler(e *core.RequestEvent) error {
	var req model.RecommendStreamRequest
	if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Channel == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "channel is required",
		})
	}

	resp, err := h.service.Stream().GetRecommendedChannels(&req)
	if err != nil {
		h.logger.Error("failed to get recommended channels", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, resp)
}

func (h *Handler) GetAllStreamHandler(e *core.RequestEvent) error {
	var req model.AllStreamsRequest
	if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// Default to page 1 if not provided
	if req.Page < 1 {
		req.Page = 1
	}

	resp, err := h.service.Stream().GetAllStreams(&req)
	if err != nil {
		h.logger.Error("failed to get all streams", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, resp)
}

func (h *Handler) GetCategoriesHandler(e *core.RequestEvent) error {
	categories, err := h.service.Stream().GetCategories()
	if err != nil {
		h.logger.Error("failed to get categories", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"categories": categories,
	})
}

func (h *Handler) GetCountriesHandler(e *core.RequestEvent) error {
	countries, err := h.service.Stream().GetCountries()
	if err != nil {
		h.logger.Error("failed to get countries", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"countries": countries,
	})
}

func (h *Handler) GetLanguagesHandler(e *core.RequestEvent) error {
	languages, err := h.service.Stream().GetLanguages()
	if err != nil {
		h.logger.Error("failed to get languages", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"languages": languages,
	})
}

func (h *Handler) SearchStreamHandler(e *core.RequestEvent) error {
	var req model.SearchStreamRequest
	if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	resp, err := h.service.Stream().SearchStreams(&req)
	if err != nil {
		h.logger.Error("failed to search streams", "error", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, resp)
}

func (h *Handler) PlayStreamHandler(e *core.RequestEvent) error {
	var req model.PlayStreamRequest
	if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// At least one of token, channel_id, or url must be provided
	if req.Token == "" && req.ChannelID == "" && req.URL == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "token, channel_id, or url is required",
		})
	}

	resp, err := h.service.Stream().PlayStream(&req)
	if err != nil {
		h.logger.Error("failed to play stream", "error", err)
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, resp)
}
