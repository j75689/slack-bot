package middleware

import (
	"fmt"
	"net/http"

	"github.com/j75689/slack-bot/appruntime"
	"github.com/j75689/slack-bot/manager"
	"github.com/gin-gonic/gin"
)

// VerifyProjectMiddleware handle verify project exists
func VerifyProjectMiddleware(management *manager.Management) gin.HandlerFunc {
	return func(c *gin.Context) {

		projectName := c.Param("project")
		_, projectManager := management.Get(manager.ProjectKind)
		if pjm, ok := projectManager.(*manager.ProjectManager); ok {
			if !pjm.VerifyProject(projectName) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": fmt.Sprintf("project [%s] not found", projectName),
				})
				c.Abort()
				return
			}
		} else {
			appruntime.Logger.Error("load project manager falid")
			c.Abort()
			return
		}
	}
}
