package team_add

import (
	"avito-test-2025/internal/pkg/api"
	"avito-test-2025/internal/repository/postgresql"
	"errors"
	"fmt"
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
		TeamName string `json:"team_name"`
		Members  []struct {
			ID       string `json:"user_id"`
			Name     string `json:"username"`
			IsActive bool   `json:"is_active"`
		} `json:"members"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	members := make([]postgresql.User, 0, len(req.Members))
	for _, m := range req.Members {
		members = append(members, postgresql.User{
			ID:       m.ID,
			Name:     m.Name,
			IsActive: m.IsActive,
		})
	}

	err := h.repository.AddTeam(c.Request.Context(), req.TeamName, members)
	if err != nil {
		if errors.Is(err, postgresql.ErrTeamExists) {
			c.JSON(http.StatusBadRequest, api.ErrorResponse{Error: api.Error{
				Code:    "TEAM_EXISTS",
				Message: fmt.Sprintf("%s already exists", req.TeamName),
			}})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Error: api.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"team": req.TeamName,
	})
}
