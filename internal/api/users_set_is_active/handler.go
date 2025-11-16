package users_set_is_active

import (
	"avito-test-2025/internal/pkg/api"
	"avito-test-2025/internal/repository/postgresql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	repository Repository
}

func New(r Repository) *Handler {
	return &Handler{
		repository: r,
	}
}

func (h *Handler) Handle(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := h.repository.SetIsActive(c.Request.Context(), req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse{Error: api.Error{
				Code:    "NOT_FOUND",
				Message: "Not found",
			}})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Error: api.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"user_id":   req.UserID,
			"username":  user.Name,
			"team_name": user.TeamName,
			"is_active": req.IsActive,
		},
	})
}
