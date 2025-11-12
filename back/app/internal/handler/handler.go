package handler

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	"gitlab.yurtal.tech/company/blitz/business-card/back/internal/config"
	"gitlab.yurtal.tech/company/blitz/business-card/back/internal/service"
)

type Handler struct {
	logger  *slog.Logger
	service service.I
	cfg     *config.Config
}

func (h *Handler) Register(router *router.Router[*core.RequestEvent]) {
	// Health check endpoint
	router.GET("/api/health", h.HealthCheckHandler)
	
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/", h.AuthHandler)
		}
		stream := api.Group("/stream")
		{
			stream.POST("/watch", h.WatchStreamHandler)
			stream.POST("/play", h.PlayStreamHandler)
			stream.GET("/featured", h.FeaturedStreamHandler)
			stream.GET("/channel/:name", h.GetChannelHandler)
			stream.POST("/recommend", h.RecommendStreamHandler)
			stream.POST("/category", h.CategoryStreamHandler)
			stream.POST("/all", h.GetAllStreamHandler)
			stream.GET("/categories", h.GetCategoriesHandler)
			stream.GET("/countries", h.GetCountriesHandler)
			stream.GET("/languages", h.GetLanguagesHandler)
			stream.POST("/search", h.SearchStreamHandler)
		}

	}
}

func NewHandler(logger *slog.Logger, service service.I, cfg *config.Config) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
		cfg:     cfg,
	}
}

func (h *Handler) HealthCheckHandler(e *core.RequestEvent) error {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "stream-api",
		"version":   "1.0.0",
	}

	data, err := json.Marshal(health)
	if err != nil {
		return e.String(500, "Internal server error")
	}

	return e.String(200, string(data))
}
