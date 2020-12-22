package model

type CreateSubscriptionRequest struct {
	RegistrationId string `json:"registration_id" binding:"required"`
}

type Subscription struct {
	RegistrationId string `json:"registration_id"`
}
