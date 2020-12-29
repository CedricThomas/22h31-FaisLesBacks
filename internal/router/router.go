package router

import (
	"context"
	"net/http"
	"time"

	"github.com/appleboy/go-fcm"
	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
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
	cfg            *config.Config
}

func NewRouter(logger *logrus.Logger, engine *gin.Engine, store *airtable.Storer, fcmClient *fcm.Client, cfg *config.Config) *Router {
	r := &Router{
		logger:         logger,
		engine:         engine,
		authMiddleware: middleware.Auth0(cfg.Certificate, cfg.Audience, cfg.Issuer),
		store:          store,
		fcmClient:      fcmClient,
		cfg:            cfg,
	}
	r.registerRoute()
	return r
}

func (r *Router) registerRoute() {
	r.registerMemoRouter()
	r.registerReminderRouter()
	r.registerSubscriptionRouter()
}

func (r *Router) RegisterProcess(g *run.Group) {
	srv := &http.Server{
		Addr:    r.cfg.Port,
		Handler: r.engine,
	}
	g.Add(func() error {
		err := srv.ListenAndServe()
		if err != nil {
			r.logger.WithError(err).Error("gin runtime error")
			return err
		}
		return nil
	}, func(err error) {
		cancellableCtx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelCtx()
		_ = srv.Shutdown(cancellableCtx)
	})
}
