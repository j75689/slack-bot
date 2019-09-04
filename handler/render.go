package handler

import (
	"github.com/j75689/slack-bot/model"
	"github.com/j75689/slack-bot/tool"
)

// RenderProcesser process render
type RenderProcesser struct {
}

// Run render stage
func (obj *RenderProcesser) Run(stage *model.Stage, variables *map[string]interface{}) (string, error) {
	return tool.ExcuteGoTemplate(stage.Template, *variables)
}

func newRenderProcesser() *RenderProcesser {
	return &RenderProcesser{}
}
