package airtable

import (
	"fmt"

	"github.com/brianloveswords/airtable"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/subscription"
)

func (at *Storer) NewSubscription(registrationId, userId string) (*subscription.Subscription, error) {
	entity := subscription.Subscription{
		Fields: subscription.Fields{
			RegistrationId: registrationId,
			UserId:         userId,
		},
	}
	table := at.client.Table(entity.TableName())
	if err := table.Create(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (at *Storer) GetSubscription(registrationId string) (*subscription.Subscription, error) {
	table := at.client.Table(subscription.Subscription{}.TableName())
	var entity subscription.Subscription
	if err := table.Get(registrationId, &entity); isNotFoundErr(err) {
		return nil, model.NoSuchEntity
	} else if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (at *Storer) ListSubscription(userId string) ([]subscription.Subscription, error) {
	table := at.client.Table(subscription.Subscription{}.TableName())
	var entities []subscription.Subscription
	if err := table.List(&entities, &airtable.Options{
		Filter: fmt.Sprintf("{user_id} = \"%s\"", userId),
	}); err != nil {
		return nil, err
	}
	return entities, nil
}

func (at *Storer) DeleteSubscription(registrationId string) error {
	table := at.client.Table(subscription.Subscription{}.TableName())
	if err := table.Delete(&subscription.Subscription{
		Record: airtable.Record{
			ID: registrationId,
		},
	}); isNotFoundErr(err) {
		return model.NoSuchEntity
	} else if err != nil {
		return err
	}
	return nil
}
