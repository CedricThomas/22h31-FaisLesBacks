package cron

import (
	"context"
	"time"

	"github.com/oklog/run"
)

func (c *Cron) registerWorker(g *run.Group, cron cronable) {
	cancellableCtx, cancelCtx := context.WithCancel(context.Background())
	g.Add(func() error {
		var loop = true
		for loop {
			cron.Trigger(cancellableCtx)
			select {
			case <-time.After(cron.GetDelay()):
			case <-cancellableCtx.Done():
				loop = false
			}
		}
		return nil
	}, func(err error) {
		cancelCtx()
	})
}
