package modules

import (
	"bytes"
	"net/http"

	"github.com/j75689/slack-bot/manager"

	"github.com/j75689/slack-bot/appruntime"

	"gopkg.in/yaml.v2"

	"github.com/j75689/slack-bot/model"

	"github.com/gin-gonic/gin"
)

// HandleDryRun debug dryrun
func HandleDryRun(management *manager.Management) func(*gin.Context) {
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
			if ok, manager := management.Get(config.Kind); ok {
				reply, _ := manager.DryRun(&config)
				c.String(http.StatusOK, "---")
				c.String(http.StatusOK, reply)
				appruntime.Logger.Debug(reply)
			}

		}
	}
}
