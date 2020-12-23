package subscription

import (
	"github.com/brianloveswords/airtable"

	"github.com/CedricThomas/22h31-FaisLesBacks/api/model"
)

type (
	Fields struct {
		RegistrationId string `json:"registration_id"`
		UserId         string `json:"user_id"`
	}
	Subscription struct {
		airtable.Record
		Fields Fields
	}
)

func (Subscription) TableName() string {
	return "subscription"
}

func (s *Subscription) ToModel() *model.Subscription {
	return &model.Subscription{
		SubscriptionId:   s.ID,
		RegistrationId:   s.Fields.RegistrationId,
		RegistrationDate: s.CreatedTime,
	}
}
