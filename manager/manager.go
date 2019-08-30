package manager

import (
	"strings"

	"github.com/j75689/slack-bot/model"
)

// Manager ...
type Manager interface {
	Register(project string, config *model.SlackBotConfig) (bool, error)
	Deregister(project, configName string) (bool, error)
	Execute(project, cmd string) (string, error)
	DryRun(config *model.SlackBotConfig) (string, error)
}

// Management manager
type Management struct {
	managers map[string]Manager
}

// Get Manager
func (obj *Management) Get(kind string) (bool, Manager) {
	kind = strings.Title(strings.ToLower(kind))
	return obj.managers[kind] != nil, obj.managers[kind]
}

// NewManagement instance
func NewManagement() *Management {
	return &Management{
		managers: map[string]Manager{
			ProjectKind:     newProjectManager(),
			MessageKind:     newMessageManager(),
			InteractiveKind: nil,
		},
	}
}
