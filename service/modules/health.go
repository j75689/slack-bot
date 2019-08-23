package modules

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleHealthCheck check service health
func HandleHealthCheck() func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health": true})
	}
}
