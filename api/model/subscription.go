package model

import "time"

type CreateSubscriptionRequest struct {
	RegistrationId string `json:"registration_id" binding:"required"`
}

type Subscription struct {
	SubscriptionId   string    `json:"subscription_id"`
	RegistrationId   string    `json:"registration_id"`
	RegistrationDate time.Time `json:"registration_date"`
}
