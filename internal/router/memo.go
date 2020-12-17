package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/pkg/middleware"
)

func (r *Router) registerMemoRouter() {
	r.engine.POST("/memo", r.authMiddleware, r.handleCreateMemo)
	r.engine.GET("/memo", r.authMiddleware, r.handleListMemo)
}

type CreateMemo struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (r *Router) handleCreateMemo(c *gin.Context) {
	var req CreateMemo
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	memo, err := r.store.NewMemo(req.Title, req.Content, c.MustGet(middleware.Subject).(string))
	if err != nil {
		r.logger.WithError(err).Error("unable to create memo in the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"memo_id": memo.ID})
}

func (r *Router) handleListMemo(c *gin.Context) {
	memos, err := r.store.ListMemo(c.MustGet(middleware.Subject).(string))
	if err != nil {
		r.logger.WithError(err).Error("unable to list memo from the store")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memos)
}
