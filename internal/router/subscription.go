package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CedricThomas/22h31-FaisLesBacks/api/model"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/pkg/middleware"
	modelstore "github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model"
)

func (r *Router) registerSubscriptionRouter() {
	r.engine.POST("/subscription", r.authMiddleware, r.handleCreateSubscription)
	r.engine.GET("/subscription", r.authMiddleware, r.handleListSubscription)
	r.engine.DELETE("/subscription/:subscriptionId", r.authMiddleware, r.handleDeleteSubscription)
}

func (r *Router) handleCreateSubscription(c *gin.Context) {
	var req model.CreateSubscriptionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger := r.logger.WithField("registration_id", req.RegistrationId)
	subs, err := r.store.ListSubscription(c.MustGet(middleware.Subject).(string))
	if err != nil {
		logger.WithError(err).Error("unable to list subscription from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, subRegistered := range subs {
		if subRegistered.Fields.RegistrationId == req.RegistrationId {
			c.JSON(http.StatusAlreadyReported, subRegistered.ToModel())
			return
		}
	}
	sub, err := r.store.NewSubscription(req.RegistrationId, c.MustGet(middleware.Subject).(string))
	if err != nil {
		logger.WithError(err).Error("unable to create subscription")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sub.ToModel())
	return
}

func (r *Router) handleListSubscription(c *gin.Context) {
	subs, err := r.store.ListSubscription(c.MustGet(middleware.Subject).(string))
	if err != nil {
		r.logger.WithError(err).Error("unable to list subscriptions from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]*model.Subscription, 0, len(subs))
	for _, sub := range subs {
		resp = append(resp, sub.ToModel())
	}
	c.JSON(http.StatusOK, resp)
}

func (r *Router) handleDeleteSubscription(c *gin.Context) {
	subscriptionId := c.Params.ByName("subscriptionId")
	logger := r.logger.WithField("subscription_id", subscriptionId)
	if len(subscriptionId) == 0 {
		logger.Error("invalid subscription id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}
	sub, err := r.store.GetSubscription(subscriptionId)
	if err == modelstore.NoSuchEntity {
		logger.WithError(err).Error("unable to find subscription in the store")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to get subscription from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sub.Fields.UserId != c.MustGet(middleware.Subject).(string) {
		c.JSON(http.StatusNotFound, gin.H{"error": modelstore.NoSuchEntity.Error()})
		return
	}
	if err := r.store.DeleteSubscription(subscriptionId); err != nil {
		logger.WithError(err).Error("unable to remove subscription in the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
