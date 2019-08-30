package manager

import (
	"github.com/j75689/slack-bot/appruntime"
	"github.com/j75689/slack-bot/handler"
	"github.com/j75689/slack-bot/model"
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
func (obj *MessageManager) Register(project string, config *model.SlackBotConfig) (ok bool, err error) {
	rollback := false
	index := obj.findIndex(project)
	for _, cmd := range config.Task.Command {
		err = index.Insert(cmd, []byte(config.MetaData.Name))
	}
	rollback = err != nil

	if !rollback {
		err = appruntime.DB.Save(project, config.Kind, config.MetaData.Name, config)
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
func (obj *MessageManager) Deregister(project, configName string) (ok bool, err error) {
	var config model.SlackBotConfig
	data, _ := appruntime.DB.Find(project, MessageKind, configName)
	if err = yaml.Unmarshal(data, &config); err != nil {
		return false, err
	}
	index := obj.findIndex(project)
	for _, cmd := range config.Task.Command {
		index.Delete(cmd)
	}

	appruntime.DB.Delete(project, MessageKind, configName)

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
	data, _ := appruntime.DB.Find(project, MessageKind, string(configName))
	if err = yaml.Unmarshal(data, &config); err != nil {
		return "", err
	}

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
