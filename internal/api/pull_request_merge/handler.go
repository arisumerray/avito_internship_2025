package pull_request_merge

import (
	"avito-test-2025/internal/pkg/api"
	"avito-test-2025/internal/repository/postgresql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
		ID string `json:"pull_request_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	pr, err := h.repository.Merge(c.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse{Error: api.Error{
				Code:    "NOT_FOUND",
				Message: "PR not found",
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
	c.JSON(http.StatusOK, gin.H{
		"pr": gin.H{
			"pull_request_id":    req.ID,
			"pull_request_name":  pr.Name,
			"author_id":          pr.AuthorID,
			"status":             "MERGED",
			"assigned_reviewers": assigned,
			"merged_at":          pr.MergedAt.Format(time.RFC3339),
		},
	})
}
