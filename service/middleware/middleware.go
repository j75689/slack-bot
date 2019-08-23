package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NoMethodHandler handle method not found
func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method Not Allowed"})
	}
}

// NoRouteHandler handle route path not found
func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page Not Found"})
	}
}
