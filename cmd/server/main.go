package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/pkg/config"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/router"
)

func main() {
	logger := newLogger()
	cfg, err := config.NewConfig()
	if err != nil {
		logger.WithError(err).Error("unable to parse config")
		os.Exit(1)
	}
	logger.WithField("config", cfg.String()).Info("configuration loaded")
	engine := gin.Default()
	r, err := router.NewRouter(logger, engine, cfg)
	if err != nil {
		logger.WithError(err).Error("cannot create router")
		os.Exit(1)
	}
	r.RegisterRoute()
	if err := engine.Run(cfg.Port); err != nil {
		logger.WithError(err).Error("runtime error")
		os.Exit(1)
	}
}
