package handler

import (
	"github.com/j75689/slack-bot/model"
)

// StageType ...
type StageType string

const (
	// stage type
	actionTag StageType = "action"
	renderTag StageType = "render"
)

// Handler handle slack message
type Handler interface {
	Do(config *model.SlackBotConfig) (string, error)
}

// StageRunner run stage action
type StageRunner interface {
	Run(stage *model.Stage, variables *map[string]interface{}) (string, error)
}
