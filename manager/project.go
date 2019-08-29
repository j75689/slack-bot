package manager

import (
	"errors"

	"github.com/Invisibi-nd/slack-bot/appruntime"
	"github.com/Invisibi-nd/slack-bot/model"
)

// ProjectManager manage project resource
type ProjectManager struct {
}

// VerifyProject check project existes
func (obj *ProjectManager) VerifyProject(project string) bool {
	return appruntime.DB.CheckProject(project)
}

// Register config
func (obj *ProjectManager) Register(project string, config *model.SlackBotConfig) (ok bool, err error) {
	err = appruntime.DB.Save(project, config.Kind, config.MetaData.Name, config)
	return err == nil, err
}

// Deregister config
func (obj *ProjectManager) Deregister(project, configName string) (ok bool, err error) {
	err = appruntime.DB.Delete(project, ProjectKind, configName)
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
