package manager

import (
	"errors"

	"github.com/j75689/slack-bot/kind"

	"github.com/j75689/slack-bot/appruntime"
	"github.com/j75689/slack-bot/model"
)

// ProjectManager manage project resource
type ProjectManager struct {
}

// VerifyProject check project existes
func (obj *ProjectManager) VerifyProject(project string) bool {
	return appruntime.DB.CheckProject(project)
}

// Register config
func (obj *ProjectManager) Register(config *model.SlackBotConfig) (ok bool, err error) {
	err = appruntime.DB.Save(config.MetaData.Name, config.Kind, config.MetaData.Name, config)
	return err == nil, err
}

// Deregister config
func (obj *ProjectManager) Deregister(config *model.SlackBotConfig) (ok bool, err error) {
	err = appruntime.DB.Delete(config.MetaData.Name, kind.Project, config.MetaData.Name)
	return err == nil, err
}

// Execute command
func (obj *ProjectManager) Execute(project, cmd string) (reply string, err error) {
	return "", errors.New("unsupported method")

}

// DryRun test config
func (obj *ProjectManager) DryRun(config *model.SlackBotConfig) (string, error) {
	return "", errors.New("unsupported method")
}

func newProjectManager() *ProjectManager {
	return &ProjectManager{}
}
