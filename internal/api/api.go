package api

import (
	"avito-test-2025/internal/api/health"
	"avito-test-2025/internal/api/pull_request_create"
	"avito-test-2025/internal/api/pull_request_merge"
	"avito-test-2025/internal/api/pull_request_reassign"
	"avito-test-2025/internal/api/team_add"
	"avito-test-2025/internal/api/team_deactivate"
	"avito-test-2025/internal/api/team_get"
	"avito-test-2025/internal/api/users_get_review"
	"avito-test-2025/internal/api/users_set_is_active"

	"github.com/gin-gonic/gin"
)

type Api struct {
	repository Repository
}

func New(repository Repository) *Api {
	return &Api{
		repository: repository,
	}
}

func (a *Api) RegisterRoutes(r *gin.Engine) {
	reassignHandler := pull_request_reassign.New(a.repository)
	createHandler := pull_request_create.New(a.repository)
	mergeHandler := pull_request_merge.New(a.repository)
	teamAddHandler := team_add.New(a.repository)
	teamGetHandler := team_get.New(a.repository)
	setActiveHandler := users_set_is_active.New(a.repository)
	getReviewHandler := users_get_review.New(a.repository)
	teamDeactivateHandler := team_deactivate.New(a.repository)

	r.POST("/team/add", teamAddHandler.Handle)
	r.GET("/team/get", teamGetHandler.Handle)
	r.POST("/team/deactivate", teamDeactivateHandler.Handle)

	r.POST("/users/setIsActive", setActiveHandler.Handle)
	r.GET("/users/getReview", getReviewHandler.Handle)

	r.POST("/pullRequest/create", createHandler.Handle)
	r.POST("/pullRequest/merge", mergeHandler.Handle)
	r.POST("/pullRequest/reassign", reassignHandler.Handle)

	r.GET("/health", health.Handle)
}
