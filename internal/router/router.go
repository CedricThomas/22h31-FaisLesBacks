package router

import (
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
}

func NewRouter(logger *logrus.Logger, engine *gin.Engine, cfg *config.Config) *Router {
	return &Router{
		logger:         logger,
		engine:         engine,
		authMiddleware: middleware.Auth0(cfg.Certificate, cfg.Audience, cfg.Issuer),
		store:          airtable.New(cfg.ApiKey, cfg.BaseID),
	}
}

func (r *Router) RegisterRoute() {
	r.registerMemoRouter()
}
