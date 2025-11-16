package pull_request_create

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
		ID   string `json:"pull_request_id"`
		Name string `json:"pull_request_name"`
		Auth string `json:"author_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	pr, err := h.repository.Create(c.Request.Context(), req.ID, req.Name, req.Auth)
	if err != nil {
		switch {
		case errors.Is(err, postgresql.ErrPRExists):
			c.JSON(http.StatusConflict, api.ErrorResponse{Error: api.Error{
				Code:    "PR_EXISTS",
				Message: "PR id already exists",
			}})
			return
		case errors.Is(err, postgresql.ErrNotFound):
			c.JSON(http.StatusNotFound, api.ErrorResponse{Error: api.Error{
				Code:    "NOT_FOUND",
				Message: "Team/author not found",
			}})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Error: api.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}})
		return
	}
	assigned := make([]string, 0, 2)
	if pr.Reviewer1ID != "" {
		assigned = append(assigned, pr.Reviewer1ID)
	}
	if pr.Reviewer2ID != "" {
		assigned = append(assigned, pr.Reviewer2ID)
	}
	c.JSON(http.StatusCreated, gin.H{
		"pr": gin.H{
			"pull_request_id":    req.ID,
			"pull_request_name":  req.Name,
			"author_id":          req.Auth,
			"status":             "OPEN",
			"assigned_reviewers": assigned,
		},
	})
}
