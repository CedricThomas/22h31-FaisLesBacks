package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CedricThomas/22h31-FaisLesBacks/api/model"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/pkg/middleware"
	modelstore "github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/memo"
)

func (r *Router) registerReminderRouter() {
	r.engine.POST("/memo/:memoId/reminder", r.authMiddleware, r.handleCreateReminder)
	r.engine.GET("/memo/:memoId/reminder", r.authMiddleware, r.handleListReminder)
}

func (r *Router) getMemoById(memoId string, userId string) (*memo.Memo, error) {
	mem, err := r.store.GetMemo(memoId)
	if err != nil {
		r.logger.WithError(err).Error("unable to memo from store")
		return nil, err
	}
	if mem.Fields.UserId != userId {
		r.logger.Error("memo found but with invalid owner")
		return nil, modelstore.NoSuchEntity
	}
	return mem, nil
}

func (r *Router) handleCreateReminder(c *gin.Context) {
	memoId := c.Params.ByName("memoId")
	logger := r.logger.WithField("memo_id", memoId)
	if len(memoId) == 0 {
		logger.Error("invalid memo id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid memo id"})
		return
	}
	var req model.CreateReminderRequest
	if err := c.BindJSON(&req); err != nil {
		logger.WithError(err).Error("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if time.Now().After(req.Date) {
		logger.Error("cannot create a reminder for the past")
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot create a reminder for the past"})
		return
	}
	if _, err := r.getMemoById(memoId, c.MustGet(middleware.Subject).(string)); err == modelstore.NoSuchEntity {
		logger.WithError(err).Error("cannot find memo")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to get memo from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reminder, err := r.store.NewReminder(memoId, req.Title, req.Content, req.Date)
	if err != nil {
		logger.WithError(err).Error("unable to create reminder")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reminder)
}

func (r *Router) handleListReminder(c *gin.Context) {
	memoId := c.Params.ByName("memoId")
	logger := r.logger.WithField("memo_id", memoId)
	if len(memoId) == 0 {
		logger.Error("invalid memo id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid memo id"})
		return
	}
	if _, err := r.getMemoById(memoId, c.MustGet(middleware.Subject).(string)); err == modelstore.NoSuchEntity {
		logger.WithError(err).Error("cannot find memo")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to get memo from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reminders, err := r.store.ListReminder(memoId)
	if err != nil {
		logger.WithError(err).Error("unable to list reminder from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]*model.Reminder, 0, len(reminders))
	for _, rem := range reminders {
		resp = append(resp, &model.Reminder{
			Id:        rem.ID,
			CreatedAt: rem.CreatedTime,
			Title:     rem.Fields.Title,
			Content:   rem.Fields.Content,
			Date:      rem.Fields.ReminderDate,
			Triggered: rem.Fields.Triggered,
		})
	}
	c.JSON(http.StatusOK, resp)
}
