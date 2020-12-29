package reminder

import (
	"context"
	"fmt"
	"time"

	"github.com/appleboy/go-fcm"
	"github.com/sirupsen/logrus"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/airtable"
)

type Reminder struct {
	logger    *logrus.Logger
	store     *airtable.Storer
	fcmClient *fcm.Client
}

func New(logger *logrus.Logger, store *airtable.Storer, fcmClient *fcm.Client) *Reminder {
	return &Reminder{
		logger:    logger,
		store:     store,
		fcmClient: fcmClient,
	}
}

func (r *Reminder) GetDelay() time.Duration {
	return time.Minute
}

func (r *Reminder) Trigger(ctx context.Context) {
	rems, err := r.store.ListReminderToTrigger()
	if err != nil {
		r.logger.WithError(err).Error("unable to list reminder from store")
		return
	}
	fmt.Println(rems)
	_ = rems
}
