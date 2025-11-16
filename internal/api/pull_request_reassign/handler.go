package pull_request_reassign

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"avito-test-2025/internal/pkg/api"
	repo "avito-test-2025/internal/repository/postgresql"
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
		ID        string `json:"pull_request_id"`
		OldUserID string `json:"old_user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	_, newUserID, err := h.repository.Reassign(c.Request.Context(), req.ID, req.OldUserID)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			c.JSON(http.StatusNotFound, api.ErrorResponse{
				Error: api.Error{
					Code:    "NOT_FOUND",
					Message: "pr or user is not found",
				},
			})
			return
		case errors.Is(err, repo.ErrNoCandidates):
			c.JSON(http.StatusNotFound, api.ErrorResponse{
				Error: api.Error{
					Code:    "NO_CANDIDATE",
					Message: "no active replacement candidate in team",
				},
			})
			return
		case errors.Is(err, repo.ErrPRIsMerged):
			c.JSON(http.StatusNotFound, api.ErrorResponse{
				Error: api.Error{
					Code:    "PR_MERGED",
					Message: "cannot reassign on merged PR",
				},
			})
			return
		case errors.Is(err, repo.ErrNotAssigned):
			c.JSON(http.StatusNotFound, api.ErrorResponse{
				Error: api.Error{
					Code:    "NOT_ASSIGNED",
					Message: "reviewer is not assigned to this PR ",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Error: api.Error{
			Code:    "INTERNAL_ERROR",
			Message: "internal error",
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"pr": gin.H{
			"pull_request_id": req.ID,
		},
		"replaced_by": newUserID,
	})
}
