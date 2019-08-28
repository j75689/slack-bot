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
	return appruntime.DB.CheckTable(project)
}

// Register config
func (obj *ProjectManager) Register(project string, config *model.SlackBotConfig) (ok bool, err error) {
	err = appruntime.DB.Save(project, config.MetaData.Name, config)
	return err == nil, err
}

// Deregister config
func (obj *ProjectManager) Deregister(project string, configName string) (ok bool, err error) {
	err = appruntime.DB.Delete(project, configName)
	return err == nil, err
}

// Execute command
func (obj *ProjectManager) Execute(project string, cmd string) (reply string, err error) {
	return "", errors.New("unsupported method")

}

// DryRun test config
func (obj *ProjectManager) DryRun(config *model.SlackBotConfig) (string, error) {
	return "", errors.New("unsupported method")
}

func newProjectManager() *ProjectManager {
	return &ProjectManager{}
}
