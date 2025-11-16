package team_get

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
	teamName := c.Query("team_name")

	team, err := h.repository.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse{Error: api.Error{
				Code:    "NOT_FOUND",
				Message: "team not found",
			}})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Error: api.Error{
			Message: err.Error(),
		}})
		return
	}
	members := make([]struct {
		ID       string `json:"user_id"`
		Name     string `json:"username"`
		IsActive bool   `json:"is_active"`
	}, 0, len(team.Members))
	for _, member := range team.Members {
		members = append(members, struct {
			ID       string `json:"user_id"`
			Name     string `json:"username"`
			IsActive bool   `json:"is_active"`
		}{
			ID:       member.ID,
			Name:     member.Name,
			IsActive: member.IsActive,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"team_name": teamName,
		"members":   members,
	})
}
