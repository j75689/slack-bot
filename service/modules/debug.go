package modules

import (
	"bytes"
	"net/http"

	"github.com/Invisibi-nd/slack-bot/appruntime"

	"gopkg.in/yaml.v2"

	"github.com/Invisibi-nd/slack-bot/model"

	"github.com/Invisibi-nd/slack-bot/handler"
	"github.com/gin-gonic/gin"
)

// HandleDryRun debug dryrun
func HandleDryRun(manager *handler.Manager) func(*gin.Context) {
	return func(c *gin.Context) {
		data, _ := c.GetRawData()
		reader := bytes.NewReader(data)
		dec := yaml.NewDecoder(reader)
		var config model.SlackBotConfig

		for {
			err := dec.Decode(&config)
			if err != nil {
				break
			}
			reply, _ := manager.DryRun(&config)
			c.String(http.StatusOK, "---")
			c.String(http.StatusOK, reply)
			appruntime.Logger.Debug(reply)

		}
	}
}
