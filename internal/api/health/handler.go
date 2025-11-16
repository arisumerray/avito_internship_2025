package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
