package handler

import (
	"github.com/Invisibi-nd/slack-bot/appruntime"
	"github.com/Invisibi-nd/slack-bot/model"
	"github.com/Invisibi-nd/slack-bot/tool/tree"
)

// Manager manage handlers
type Manager struct {
	Index    *tree.Tree
	Handlers []Handler
}

// Register config
func (obj *Manager) Register(config *model.HandlerConfig) (ok bool, err error) {
	for _, cmd := range config.Command {
		err = obj.Index.Insert(cmd, []byte(config.ConfigID))
	}
	if err != nil {
		appruntime.Logger.Error("register config error: " + err.Error())
		// rollback
		for _, cmd := range config.Command {
			obj.Index.Delete(cmd)
		}
		return false, err
	}
	return true, nil
}

// Deregister config
func (obj *Manager) Deregister(configID string) (ok bool, err error) {
	return
}

// Execute command
func (obj *Manager) Execute(cmd string) (reply string, err error) {
	return
}

// NewManager create new handler manager
func NewManager() *Manager {
	return &Manager{
		Index: tree.NewTree(),
		Handlers: []Handler{
			newSlackMessageHandler(),
		},
	}
}
