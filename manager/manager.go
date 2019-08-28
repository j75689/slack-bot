package manager

import (
	"strings"

	"github.com/Invisibi-nd/slack-bot/model"
)

// Manager ...
type Manager interface {
	Register(project string, config *model.SlackBotConfig) (bool, error)
	Deregister(project string, configName string) (bool, error)
	Execute(project string, cmd string) (string, error)
	DryRun(config *model.SlackBotConfig) (string, error)
}

// Management manager
type Management struct {
	managers map[string]Manager
}

// Get Manager
func (obj *Management) Get(kind string) (bool, Manager) {
	kind = strings.Title(kind)
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
