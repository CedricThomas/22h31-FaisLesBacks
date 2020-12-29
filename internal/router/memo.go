package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CedricThomas/22h31-FaisLesBacks/api/model"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/pkg/middleware"
	modelstore "github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/memo"
)

func (r *Router) registerMemoRouter() {
	r.engine.POST("/memo", r.authMiddleware, r.handleCreateMemo)
	r.engine.GET("/memo", r.authMiddleware, r.handleListMemo)
	r.engine.GET("/memo/:memoId", r.authMiddleware, r.handleGetMemo)
	r.engine.PUT("/memo/:memoId", r.authMiddleware, r.handleUpdateMemo)
	r.engine.DELETE("/memo/:memoId", r.authMiddleware, r.handleDeleteMemo)
}

func (r *Router) handleCreateMemo(c *gin.Context) {
	var req model.CreateMemoRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mem, err := r.store.NewMemo(req.Title, req.Content, c.MustGet(middleware.Subject).(string))
	if err != nil {
		r.logger.WithError(err).Error("unable to create memo in the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mem.ToModel())
}

func (r *Router) handleListMemo(c *gin.Context) {
	memos, err := r.store.ListMemo(c.MustGet(middleware.Subject).(string))
	if err != nil {
		r.logger.WithError(err).Error("unable to list memo from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]*model.Memo, 0, len(memos))
	for _, mem := range memos {
		resp = append(resp, mem.ToModel())
	}
	c.JSON(http.StatusOK, resp)
}

func (r *Router) handleGetMemo(c *gin.Context) {
	memoId := c.Params.ByName("memoId")
	logger := r.logger.WithField("memo_id", memoId)
	if len(memoId) == 0 {
		logger.Error("invalid memo id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid memo id"})
		return
	}
	mem, err := r.store.GetMemo(memoId)
	if err == modelstore.NoSuchEntity {
		logger.WithError(err).Error("unable to find memo in the store")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to get memo from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if mem.Fields.UserId != c.MustGet(middleware.Subject).(string) {
		c.JSON(http.StatusNotFound, gin.H{"error": modelstore.NoSuchEntity.Error()})
		return
	}
	c.JSON(http.StatusOK, mem.ToModel())
}

func (r *Router) handleUpdateMemo(c *gin.Context) {
	memoId := c.Params.ByName("memoId")
	logger := r.logger.WithField("memo_id", memoId)
	if len(memoId) == 0 {
		logger.Error("invalid memo id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid memo id"})
		return
	}
	mem, err := r.store.GetMemo(memoId)
	if err == modelstore.NoSuchEntity {
		logger.WithError(err).Error("unable to find memo in the store")
		c.JSON(http.StatusNotFound, gin.H{"error": modelstore.NoSuchEntity.Error()})
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to get memo from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if mem.Fields.UserId != c.MustGet(middleware.Subject).(string) {
		c.JSON(http.StatusNotFound, gin.H{"error": modelstore.NoSuchEntity.Error()})
		return
	}
	var req model.UpdateMemoRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedMemo, err := r.store.UpdateMemo(memoId, &memo.Fields{
		Title:   req.Title,
		Content: req.Content,
		UserId:  c.MustGet(middleware.Subject).(string),
	})
	if err != nil {
		logger.WithError(err).Error("unable to update the memo from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedMemo.ToModel())
	return
}

func (r *Router) handleDeleteMemo(c *gin.Context) {
	memoId := c.Params.ByName("memoId")
	logger := r.logger.WithField("memo_id", memoId)
	if len(memoId) == 0 {
		logger.Error("invalid memo id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid memo id"})
		return
	}
	mem, err := r.store.GetMemo(memoId)
	if err == modelstore.NoSuchEntity {
		logger.WithError(err).Error("unable to find memo in the store")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to get memo from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if mem.Fields.UserId != c.MustGet(middleware.Subject).(string) {
		c.JSON(http.StatusNotFound, gin.H{"error": modelstore.NoSuchEntity.Error()})
		return
	}
	if err := r.store.DeleteMemo(memoId); err != nil {
		logger.WithError(err).Error("unable to remove memo in the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
