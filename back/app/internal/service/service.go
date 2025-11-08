package service

import (
	"github.com/pocketbase/pocketbase"
	"gitlab.yurtal.tech/company/blitz/business-card/back/internal/model"
)

type AuthorizationI interface {
	AmoCRMTokenExchange(req *model.AmoCRMTokenExchangeRequest) (*model.AmoCRMTokenExchangeResponse, error)
}

type StreamI interface {
	WatchStream(req *model.WatchStreamRequest) (*model.WatchStreamResponse, error)
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
	return &service{
		AuthorizationI: NewAuthorizationS(app),
		StreamI:        NewStream(app),
	}
}
