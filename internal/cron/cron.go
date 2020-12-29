package cron

import (
	"context"
	"time"

	"github.com/appleboy/go-fcm"
	"github.com/oklog/run"
	"github.com/sirupsen/logrus"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/cron/reminder"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/airtable"
)

type (
	cronable interface {
		GetDelay() time.Duration
		Trigger(ctx context.Context)
	}

	Cron struct {
		logger    *logrus.Logger
		store     *airtable.Storer
		fcmClient *fcm.Client
		crons     []cronable
	}
)

func New(logger *logrus.Logger, store *airtable.Storer, fcmClient *fcm.Client) *Cron {
	c := &Cron{
		logger:    logger,
		store:     store,
		fcmClient: fcmClient,
	}
	c.registerCrons()
	return c
}

func (c *Cron) registerCrons() {
	reminderCron := reminder.New(c.logger, c.store, c.fcmClient)
	c.crons = append(c.crons, reminderCron)
}

func (c *Cron) RegisterProcess(g *run.Group) {
	for _, cron := range c.crons {
		c.registerWorker(g, cron)
	}
}
