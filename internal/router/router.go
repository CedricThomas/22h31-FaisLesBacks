package router

import (
	"github.com/appleboy/go-fcm"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/pkg/config"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/pkg/middleware"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/airtable"
)

type Router struct {
	logger         *logrus.Logger
	engine         *gin.Engine
	authMiddleware gin.HandlerFunc
	store          store.Store
	fcmClient      *fcm.Client
}

func NewRouter(logger *logrus.Logger, engine *gin.Engine, cfg *config.Config) (*Router, error) {
	client, err := fcm.NewClient(cfg.FcmServerKey)
	if err != nil {
		return nil, err
	}
	return &Router{
		logger:         logger,
		engine:         engine,
		authMiddleware: middleware.Auth0(cfg.Certificate, cfg.Audience, cfg.Issuer),
		store:          airtable.New(cfg.ApiKey, cfg.BaseID),
		fcmClient:      client,
	}, nil
}

func (r *Router) RegisterRoute() {
	r.registerMemoRouter()
	r.registerSubscriptionRouter()
}
