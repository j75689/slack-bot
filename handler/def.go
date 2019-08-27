package handler

import (
	"github.com/Invisibi-nd/slack-bot/model"
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
	Do(config *model.HandlerConfig) (string, error)
}

// StageRunner run stage action
type StageRunner interface {
	Run(stage *model.HandlerStageConfig, variables *map[string]interface{}) (string, error)
}
