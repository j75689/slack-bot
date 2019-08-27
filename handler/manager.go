package handler

import (
	"github.com/Invisibi-nd/slack-bot/appruntime"
	"github.com/Invisibi-nd/slack-bot/model"
	"github.com/Invisibi-nd/slack-bot/tool/tree"
	"gopkg.in/yaml.v2"
)

// Manager manage handlers
type Manager struct {
	Index          map[string]*tree.Tree
	MessageHandler Handler
}

func (obj *Manager) findIndex(project string) (index *tree.Tree) {
	index = obj.Index[project]
	if index == nil {
		index = tree.NewTree()
		obj.Index[project] = index
	}
	return index
}

// Register config
func (obj *Manager) Register(project string, config *model.HandlerConfig) (ok bool, err error) {
	rollback := false
	index := obj.findIndex(project)
	for _, cmd := range config.Command {
		err = index.Insert(cmd, []byte(config.ConfigID))
	}
	rollback = err != nil

	if !rollback {
		err = appruntime.DB.Save(project, config.ConfigID, config)
		rollback = err != nil
	}

	if rollback {
		for _, cmd := range config.Command {
			index.Delete(cmd)
		}
		return
	}

	return true, nil
}

// Deregister config
func (obj *Manager) Deregister(project string, configID string) (ok bool, err error) {
	var config model.HandlerConfig
	data, _ := appruntime.DB.Find(project, configID)
	if err = yaml.Unmarshal(data, &config); err != nil {
		return false, err
	}
	index := obj.findIndex(project)
	for _, cmd := range config.Command {
		index.Delete(cmd)
	}

	appruntime.DB.Delete(project, configID)

	return true, nil
}

// Execute command
func (obj *Manager) Execute(project string, cmd string) (reply string, err error) {
	index := obj.findIndex(project)
	configID, err := index.Search(cmd)
	if err != nil {
		return "", err
	}

	var config model.HandlerConfig
	data, _ := appruntime.DB.Find(project, string(configID))
	if err = yaml.Unmarshal(data, &config); err != nil {
		return "", err
	}

	return obj.MessageHandler.Do(&config)

}

// DryRun test config
func (obj *Manager) DryRun(config *model.HandlerConfig) (string, error) {
	return obj.MessageHandler.Do(config)
}

// NewManager create new handler manager
func NewManager() *Manager {
	return &Manager{
		Index:          make(map[string]*tree.Tree),
		MessageHandler: newSlackMessageHandler(),
	}
}
