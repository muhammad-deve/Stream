package service

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"gitlab.yurtal.tech/company/blitz/business-card/back/internal/config"
	"gitlab.yurtal.tech/company/blitz/business-card/back/internal/model"
	redisClient "gitlab.yurtal.tech/company/blitz/business-card/back/internal/redis"
)

type AuthorizationI interface {
	AmoCRMTokenExchange(req *model.AmoCRMTokenExchangeRequest) (*model.AmoCRMTokenExchangeResponse, error)
}

type StreamI interface {
	WatchStream(req *model.WatchStreamRequest) (*model.WatchStreamResponse, error)
	GetFeaturedChannels() ([]*model.WatchStreamResponse, error)
	GetChannelByName(channelName string) (*model.WatchStreamResponse, error)
	GetChannelsByCategory(categoryName string) ([]*model.WatchStreamResponse, error)
	GetRecommendedChannels(req *model.RecommendStreamRequest) ([]*model.WatchStreamResponse, error)
	GetAllStreams(req *model.AllStreamsRequest) (*model.AllStreamsResponse, error)
	GetCategories() ([]string, error)
	GetCountries() ([]string, error)
	GetLanguages() ([]string, error)
	SearchStreams(req *model.SearchStreamRequest) (*model.SearchStreamResponse, error)
	PlayStream(req *model.PlayStreamRequest) (*model.PlayStreamResponse, error)
}

type I interface {
	Authorization() AuthorizationI
	Stream() StreamI
}

type service struct {
	AuthorizationI
	StreamI
}

func (s *service) Authorization() AuthorizationI {
	return s.AuthorizationI
}

func (s *service) Stream() StreamI {
	return s.StreamI
}

func NewService(app *pocketbase.PocketBase) I {
	// Initialize Redis client
	cfg := config.GetConfig()
	redis, err := redisClient.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}

	return &service{
		AuthorizationI: NewAuthorizationS(app),
		StreamI:        NewStream(app, redis),
	}
}
