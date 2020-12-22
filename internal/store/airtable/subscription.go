package airtable

import (
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/subscription"
)

func (at *Storer) NewSubscription(registrationId, userId string) (*subscription.Subscription, error) {
	entity := subscription.Subscription{
		Fields: subscription.Fields{
			RegistrationId: registrationId,
			UserId:  userId,
		},
	}
	table := at.client.Table(entity.TableName())
	if err := table.Create(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}
