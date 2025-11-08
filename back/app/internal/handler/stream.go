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
