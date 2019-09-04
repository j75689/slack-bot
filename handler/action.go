package handler

import (
	"github.com/j75689/slack-bot/appruntime"
	"github.com/j75689/slack-bot/model"
	"github.com/j75689/slack-bot/plugin"
	"github.com/j75689/slack-bot/tool"
	"gopkg.in/yaml.v2"
)

// ActionProcesser process plugin action
type ActionProcesser struct {
	Plugins *plugin.Pool
}

// Run action stage
func (obj *ActionProcesser) Run(stage *model.Stage, variables *map[string]interface{}) (string, error) {
	paramter, _ := yaml.Marshal(stage.Paramter)
	paramter = tool.ReplaceVariables(paramter, *variables)
	return "", obj.Plugins.Execute(stage.Plugin, paramter, variables, func(data interface{}) {
		(*variables)[stage.Output] = data
	})
}

func newActionProcesser() *ActionProcesser {
	return &ActionProcesser{
		Plugins: plugin.Load(appruntime.Env.PluginPath),
	}
}
