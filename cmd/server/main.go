package main

import (
	"os"

	"github.com/appleboy/go-fcm"
	"github.com/codingsince1985/geo-golang/frenchapigouv"
	"github.com/gin-gonic/gin"
	"github.com/oklog/run"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/cron"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/pkg/config"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/router"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/airtable"
)

func main() {
	logger := newLogger()
	cfg, err := config.NewConfig()
	if err != nil {
		logger.WithError(err).Error("unable to parse config")
		os.Exit(1)
	}
	logger.WithField("config", cfg.String()).Info("configuration loaded")

	geoSolver := frenchapigouv.Geocoder()
	fcmClient, err := fcm.NewClient(cfg.FcmServerKey)
	if err != nil {
		logger.WithError(err).Error("unable to create airtable store")
		os.Exit(1)
	}
	store := airtable.New(cfg.ApiKey, cfg.BaseID)

	engine := gin.Default()
	r := router.NewRouter(logger, engine, store, fcmClient, geoSolver, cfg)

	crons := cron.New(logger, store, fcmClient)

	var g run.Group
	r.RegisterProcess(&g)
	crons.RegisterProcess(&g)

	if err := g.Run(); err != nil {
		logger.WithError(err).Error("runtime error")
		os.Exit(1)
	}
}
