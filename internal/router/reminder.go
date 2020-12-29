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
	r.engine.GET("/memo/:memoId/reminder/:reminderId", r.authMiddleware, r.handleGetReminder)
	r.engine.PUT("/memo/:memoId/reminder/:reminderId", r.authMiddleware, r.handleUpdateReminder)
	r.engine.DELETE("/memo/:memoId/reminder/:reminderId", r.authMiddleware, r.handleDeleteReminder)
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
	c.JSON(http.StatusOK, reminder.ToModel())
}

func (r *Router) handleGetReminder(c *gin.Context) {
	memoId := c.Params.ByName("memoId")
	reminderId := c.Params.ByName("reminderId")
	logger := r.logger.WithField("memo_id", memoId).WithField("reminder_id", reminderId)
	if len(memoId) == 0 || len(reminderId) == 0 {
		logger.Error("invalid id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
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
	if rem, err := r.store.GetReminder(reminderId); err == modelstore.NoSuchEntity {
		logger.WithError(err).Error("cannot find reminder in the store")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to get reminder from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if rem.Fields.MemoId != memoId {
		logger.WithError(err).Error("reminder found but with invalid memo id")
		c.JSON(http.StatusNotFound, gin.H{"error": modelstore.NoSuchEntity.Error()})
		return
	}
	rem, err := r.store.GetReminder(reminderId)
	if err == modelstore.NoSuchEntity {
		logger.WithError(err).Error("cannot find reminder")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to find reminder")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rem.ToModel())
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

func (r *Router) handleUpdateReminder(c *gin.Context) {
	memoId := c.Params.ByName("memoId")
	reminderId := c.Params.ByName("reminderId")
	logger := r.logger.WithField("memo_id", memoId).WithField("reminder_id", reminderId)
	if len(memoId) == 0 || len(reminderId) == 0 {
		logger.Error("invalid id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req model.UpdateReminderRequest
	if err := c.BindJSON(&req); err != nil {
		logger.WithError(err).Error("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	rem, err := r.store.GetReminder(reminderId)
	if err == modelstore.NoSuchEntity {
		logger.WithError(err).Error("cannot find reminder in the store")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to get reminder from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if rem.Fields.MemoId != memoId {
		logger.WithError(err).Error("reminder found but with invalid memo id")
		c.JSON(http.StatusNotFound, gin.H{"error": modelstore.NoSuchEntity.Error()})
		return
	}
	rem.Fields.Title = req.Title
	rem.Fields.Content = req.Content
	rem.Fields.ReminderDate = req.Date
	updatedRem, err := r.store.UpdateReminder(reminderId, &rem.Fields)
	if err != nil {
		logger.WithError(err).Error("unable to update reminder in the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedRem.ToModel())
}

func (r *Router) handleDeleteReminder(c *gin.Context) {
	memoId := c.Params.ByName("memoId")
	reminderId := c.Params.ByName("reminderId")
	logger := r.logger.WithField("memo_id", memoId).WithField("reminder_id", reminderId)
	if len(memoId) == 0 || len(reminderId) == 0 {
		logger.Error("invalid id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
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
	if rem, err := r.store.GetReminder(reminderId); err == modelstore.NoSuchEntity {
		logger.WithError(err).Error("cannot find reminder in the store")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to get reminder from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if rem.Fields.MemoId != memoId {
		logger.WithError(err).Error("reminder found but with invalid memo id")
		c.JSON(http.StatusNotFound, gin.H{"error": modelstore.NoSuchEntity.Error()})
		return
	}
	if err := r.store.DeleteReminder(reminderId); err != nil {
		logger.WithError(err).Error("unable to remove reminder in the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
