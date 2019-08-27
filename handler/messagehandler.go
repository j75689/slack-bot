package handler

import "github.com/Invisibi-nd/slack-bot/model"

// SlackMessageHandler process Slack message
type SlackMessageHandler struct {
	Processer map[StageType]StageRunner
}

// Do ...
func (obj *SlackMessageHandler) Do(config *model.HandlerConfig) (reply string, err error) {
	for _, stage := range config.Stage {
		if runner := obj.Processer[StageType(stage.Type)]; runner != nil {
			reply, err = runner.Run(stage, &config.Variables)
		}
		if err != nil {
			break
		}
	}
	return
}

func newSlackMessageHandler() *SlackMessageHandler {
	return &SlackMessageHandler{
		Processer: map[StageType]StageRunner{
			renderTag: &RenderProcesser{},
			actionTag: &ActionProcesser{},
		},
	}
}
