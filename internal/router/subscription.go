package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CedricThomas/22h31-FaisLesBacks/api/model"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/pkg/middleware"
)

func (r *Router) registerSubscriptionRouter() {
	r.engine.POST("/subscription", r.handleCreateSubscription)
}

func (r *Router) handleCreateSubscription(c *gin.Context) {
	var req model.CreateSubscriptionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub, err := r.store.NewSubscription(req.RegistrationId, c.MustGet(middleware.Subject).(string))
	if err != nil {
		r.logger.WithError(err).Error("unable to create subscription")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sub.ToModel())
	return
}
