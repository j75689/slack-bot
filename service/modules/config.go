package modules

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/j75689/slack-bot/appruntime"

	"github.com/gin-gonic/gin"
	"github.com/j75689/slack-bot/manager"
	"github.com/j75689/slack-bot/model"
	"gopkg.in/yaml.v2"
)

// HandleApplyConfig apply config
func HandleApplyConfig(management *manager.Management) func(*gin.Context) {
	return func(c *gin.Context) {
		data, _ := c.GetRawData()
		reader := bytes.NewReader(data)
		dec := yaml.NewDecoder(reader)
		var (
			config model.SlackBotConfig
			err    error
		)
		for {
			eof := dec.Decode(&config)
			if eof != nil {
				break
			}
			if ok, manager := management.Get(config.Kind); ok {
				// reomve old cmd
				_, err = manager.Deregister(&config)
				if err != nil {
					appruntime.Logger.Error("deregister error: " + err.Error())
				}
				// register new cmd
				_, err = manager.Register(&config)
				if err != nil {
					break
				}
			} else {
				err = errors.New("invalid config")
				break
			}

		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})

	}
}

// HandleDeleteConfig delete config
func HandleDeleteConfig(management *manager.Management) func(c *gin.Context) {
	return func(c *gin.Context) {
		data, _ := c.GetRawData()
		reader := bytes.NewReader(data)
		dec := yaml.NewDecoder(reader)
		var (
			config model.SlackBotConfig
			err    error
		)
		for {
			eof := dec.Decode(&config)
			if eof != nil {
				break
			}
			if ok, manager := management.Get(config.Kind); ok {
				_, err = manager.Deregister(&config)
				if err != nil {
					break
				}
			} else {
				err = errors.New("invalid config")
				break
			}

		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
