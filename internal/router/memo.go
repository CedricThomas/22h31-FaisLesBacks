package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CedricThomas/22h31-FaisLesBacks/api/model"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/pkg/middleware"
	modelstore "github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model"
)

func (r *Router) registerMemoRouter() {
	r.engine.POST("/memo", r.authMiddleware, r.handleCreateMemo)
	r.engine.GET("/memo", r.authMiddleware, r.handleListMemo)
	r.engine.GET("/memo/:memoId", r.authMiddleware, r.handleGetMemo)
	r.engine.PUT("/memo/:memoId", r.authMiddleware, r.handleUpdateMemo)
	r.engine.DELETE("/memo/:memoId", r.authMiddleware, r.handleDeleteMemo)
}

func (r *Router) extractMemoLocation(location *model.Location) string {
	if location == nil {
		return ""
	}
	logger := r.logger.WithField("latitude", location.Latitude).WithField("longitude", location.Longitude)
	addr, err := r.geoSolver.ReverseGeocode(location.Latitude, location.Longitude)
	if err != nil {
		logger.WithError(err).Error("unable to resolve geo location")
		return "Outside france"
	}
	return fmt.Sprintf("%s %s, %s", addr.HouseNumber, addr.Street, addr.City)
}

func (r *Router) handleCreateMemo(c *gin.Context) {
	var req model.CreateMemoRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	location := r.extractMemoLocation(req.Location)
	mem, err := r.store.NewMemo(req.Title, req.Content, location, c.MustGet(middleware.Subject).(string))
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
	mem.Fields.Title = req.Title
	mem.Fields.Content = req.Content
	mem.Fields.Location = r.extractMemoLocation(req.Location)
	updatedMemo, err := r.store.UpdateMemo(mem)
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
	if err := r.store.DeleteAllReminder(memoId); err == modelstore.NoSuchEntity {
		logger.WithError(err).Error("unable to delete all reminder associated with memo")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to delete all reminder associated with memo")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := r.store.DeleteMemo(memoId); err != nil {
		logger.WithError(err).Error("unable to remove memo in the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
