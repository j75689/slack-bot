package handler

import (
	"github.com/Invisibi-nd/slack-bot/model"
)

// Manager manage handlers
type Manager struct {
	Handlers []Handler
}

// Register config
func (obj *Manager) Register(configID string, config *model.HandlerConfig) (ok bool, err error) {
	return
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
		Handlers: []Handler{
			newSlackMessageHandler(),
		},
	}
}
