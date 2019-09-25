package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/j75689/slack-bot/kind"

	"github.com/j75689/slack-bot/appruntime"
	"github.com/j75689/slack-bot/manager"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"

	"github.com/gin-gonic/gin"
)

const slackPostMessageAPI = "https://slack.com/api/chat.postMessage"

// PostSlackMessage slack post api
// https://api.slack.com/methods/chat.postMessage
func PostSlackMessage(channel string, message slack.Msg) (result map[string]interface{}, err error) {

	message.Channel = channel
	postbody, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, slackPostMessageAPI, bytes.NewBuffer(postbody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+appruntime.Env.SlackBotOauthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(response, &result)

	if !reflect.DeepEqual(result["ok"], true) {
		return result, fmt.Errorf("%v", result["error"])
	}
	return result, nil
}

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
					_, messageManager := management.Get(kind.Message)
					replyStr, err := messageManager.Execute(projectName, cmd)
					if err != nil {
						appruntime.Logger.Error(fmt.Sprintf("[slack] process cmd [%s] error : %v", cmd, err))
						return
					}
					appruntime.Logger.Debug(fmt.Sprintf("[slack] reply message:\n%v", replyStr))
					var (
						reply slack.Msg
					)
					json.Unmarshal([]byte(replyStr), &reply)
					if _, err := PostSlackMessage(ev.Channel, reply); err != nil {
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
