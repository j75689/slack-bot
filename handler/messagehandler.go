package handler

import (
	"github.com/j75689/slack-bot/model"
)

// SlackMessageHandler process Slack message
type SlackMessageHandler struct {
	Processer map[StageType]StageRunner
}

// Do ...
func (obj *SlackMessageHandler) Do(config *model.SlackBotConfig) (reply string, err error) {
	for _, stage := range config.Task.Stage {
		if runner := obj.Processer[StageType(stage.Type)]; runner != nil {
			reply, err = runner.Run(stage, &config.Task.Variables)
		}
		if err != nil {
			break
		}
	}
	return
}

// NewSlackMessageHandler process slack message
func NewSlackMessageHandler() *SlackMessageHandler {
	return &SlackMessageHandler{
		Processer: map[StageType]StageRunner{
			renderTag: &RenderProcesser{},
			actionTag: newActionProcesser(),
		},
	}
}
