package handler

import "github.com/Invisibi-nd/slack-bot/model"

// ActionProcesser process plugin action
type ActionProcesser struct {
}

// Run action stage
func (obj *ActionProcesser) Run(stage *model.Stage, variables *map[string]interface{}) (string, error) {
	return "", nil
}
