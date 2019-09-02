package handler

import (
	"github.com/j75689/slack-bot/appruntime"
	"github.com/j75689/slack-bot/model"
	"github.com/j75689/slack-bot/plugin"
)

// ActionProcesser process plugin action
type ActionProcesser struct {
	Plugins *plugin.Pool
}

// Run action stage
func (obj *ActionProcesser) Run(stage *model.Stage, variables *map[string]interface{}) (string, error) {
	return "", nil
}

func newActionProcesser() *ActionProcesser {
	return &ActionProcesser{
		Plugins: plugin.Load(appruntime.Env.PluginPath),
	}
}
