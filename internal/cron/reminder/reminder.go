package reminder

import (
	"context"
	"fmt"
	"time"

	"github.com/appleboy/go-fcm"
	"github.com/sirupsen/logrus"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/airtable"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/reminder"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/subscription"
)

type Reminder struct {
	logger            *logrus.Logger
	store             *airtable.Storer
	fcmClient         *fcm.Client
	subscriptionCache map[string][]subscription.Subscription
}

func New(logger *logrus.Logger, store *airtable.Storer, fcmClient *fcm.Client) *Reminder {
	return &Reminder{
		logger:            logger,
		store:             store,
		fcmClient:         fcmClient,
		subscriptionCache: make(map[string][]subscription.Subscription),
	}
}

func (r *Reminder) GetDelay() time.Duration {
	return time.Minute
}

func (r *Reminder) Trigger(ctx context.Context) {
	r.logger.Info("reminder crontab triggered")
	rems, err := r.store.ListReminderToTrigger()
	fmt.Println(rems, err)
	if err != nil {
		r.logger.WithError(err).Error("unable to list reminder from store")
		return
	}
	for _, rem := range rems {
		r.sendNotification(ctx, &rem)
	}
	r.subscriptionCache = make(map[string][]subscription.Subscription)
}

func (r *Reminder) getSubscriptionsWithCacheFirst(userId string) ([]subscription.Subscription, error) {
	subscriptions, inCache := r.subscriptionCache[userId]
	if inCache {
		return subscriptions, nil
	}
	subscriptions, err := r.store.ListSubscription(userId)
	if err != nil {

		return nil, err
	}
	r.subscriptionCache[userId] = subscriptions
	return subscriptions, nil
}

func (r *Reminder) sendNotification(ctx context.Context, rem *reminder.Reminder) {
	logger := r.logger.WithField("user_id", rem.Fields.UserId).WithField("user_id", rem.ID)
	notif := fcm.Notification{
		Title: rem.Fields.Title,
		Body:  rem.Fields.Content,
	}
	subscriptions, err := r.getSubscriptionsWithCacheFirst(rem.Fields.UserId)
	if err != nil {
		logger.WithError(err).Error("unable to list user subscriptions")
		return
	}
	for _, sub := range subscriptions {
		if sub.Fields.UserId != rem.Fields.UserId {
			logger.WithField("subscription_id", sub.ID).Error("mismatch user id between subscription and reminder")
			continue
		}
		if _, err := r.fcmClient.SendWithContext(ctx, &fcm.Message{
			To:           sub.Fields.RegistrationId,
			Notification: &notif,
		}); err != nil {
			logger.WithError(err).Error("unable to send message")
			continue
		}
	}
	rem.Fields.Triggered = true
	if _, err := r.store.UpdateReminder(rem.ID, &rem.Fields); err != nil {
		logger.WithError(err).Error("unable to update reminder")
	}
}
