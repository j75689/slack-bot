package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Invisibi-nd/slack-bot/appruntime"
	"github.com/Invisibi-nd/slack-bot/manager"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"

	"github.com/gin-gonic/gin"
)

// HandleSlackEvent check service challenge
func HandleSlackEvent(api *slack.Client, botID string, management *manager.Management) func(*gin.Context) {
	return func(c *gin.Context) {

		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		body := buf.String()
		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: appruntime.Env.SlackVerificationToken}))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			c.String(http.StatusOK, r.Challenge)
		}

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.MessageEvent:
				tag := fmt.Sprintf("<@%s>", botID)
				if strings.HasPrefix(ev.Text, tag) {
					cmd := strings.TrimSpace(strings.ReplaceAll(ev.Text, tag, ""))

					// process cmd
					projectName := c.Param("project")
					_, messageManager := management.Get(manager.MessageKind)
					messageManager.Execute(projectName, cmd)

					if _, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello "+cmd, true)); err != nil {
						appruntime.Logger.Error(err.Error())
					}

				}
			}
		}
		c.Next()
	}
}

// HandleSlackInteractive handle slack Interactive
func HandleSlackInteractive(management *manager.Management) func(*gin.Context) {
	return func(c *gin.Context) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		body := buf.String()
		appruntime.Logger.Info(body)
	}
}
