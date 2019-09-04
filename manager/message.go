package manager

import (
	"github.com/j75689/slack-bot/appruntime"
	"github.com/j75689/slack-bot/handler"
	"github.com/j75689/slack-bot/kind"
	"github.com/j75689/slack-bot/model"
	"github.com/j75689/slack-bot/tool"
	"github.com/j75689/slack-bot/tool/tree"
	"gopkg.in/yaml.v2"
)

// MessageManager manage message resource
type MessageManager struct {
	Index          map[string]*tree.Tree
	MessageHandler handler.Handler
}

func (obj *MessageManager) findIndex(project string) (index *tree.Tree) {
	index = obj.Index[project]
	if index == nil {
		index = tree.NewTree()
		obj.Index[project] = index
	}
	return index
}

// Register config
func (obj *MessageManager) Register(config *model.SlackBotConfig) (ok bool, err error) {
	rollback := false
	index := obj.findIndex(config.MetaData.Project)
	for _, cmd := range config.Task.Command {
		err = index.Insert(cmd, []byte(config.MetaData.Name))
	}
	rollback = err != nil

	if !rollback {
		err = appruntime.DB.Save(config.MetaData.Project, config.Kind, config.MetaData.Name, config)
		if err != nil {
			appruntime.Logger.Error(err.Error())
		}

		rollback = err != nil
	}

	if rollback {
		for _, cmd := range config.Task.Command {
			index.Delete(cmd)
		}
		return
	}

	return true, nil
}

// Deregister config
func (obj *MessageManager) Deregister(config *model.SlackBotConfig) (ok bool, err error) {
	var oldconfig model.SlackBotConfig
	data, _ := appruntime.DB.Find(config.MetaData.Project, kind.Message, config.MetaData.Name)
	if err = yaml.Unmarshal(data, &oldconfig); err != nil || data == nil {
		return false, err
	}

	index := obj.findIndex(config.MetaData.Project)
	for _, cmd := range oldconfig.Task.Command {
		index.Delete(cmd)
	}

	appruntime.DB.Delete(config.MetaData.Project, kind.Message, config.MetaData.Name)

	return true, nil
}

// Execute command
func (obj *MessageManager) Execute(project, cmd string) (reply string, err error) {
	index := obj.findIndex(project)
	configName, err := index.Search(cmd)
	if err != nil {
		return "", err
	}

	var config model.SlackBotConfig
	data, err := appruntime.DB.Find(project, kind.Message, string(configName))

	if err != nil {
		return "", err
	}

	if err = yaml.Unmarshal(data, &config); err != nil {
		return "", err
	}

	tool.ResolveVariables(cmd, config.Task.Command, &config.Task.Variables)

	return obj.MessageHandler.Do(&config)

}

// DryRun test config
func (obj *MessageManager) DryRun(config *model.SlackBotConfig) (string, error) {
	return obj.MessageHandler.Do(config)
}

func newMessageManager() *MessageManager {
	return &MessageManager{
		Index:          make(map[string]*tree.Tree),
		MessageHandler: handler.NewSlackMessageHandler(),
	}
}
