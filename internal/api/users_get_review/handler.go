package users_get_review

import (
	"avito-test-2025/internal/pkg/api"
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
	userID := c.Query("user_id")
	_ = userID

	prs, err := h.repository.GetByReviewerID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Error: api.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}})
		return
	}
	type prShort struct {
		ID       string `json:"pull_request_id"`
		Name     string `json:"pull_request_name"`
		AuthorID string `json:"author_id"`
		Status   string `json:"status"`
	}
	prShorts := make([]prShort, 0, len(prs))
	for _, pr := range prs {
		prShorts = append(prShorts, prShort{
			ID:       pr.ID,
			Name:     pr.Name,
			AuthorID: pr.AuthorID,
			Status:   pr.Status,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"pull_requests": prShorts,
	})
}
