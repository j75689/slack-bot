package handler

import (
	"bytes"
	"text/template"

	"github.com/j75689/slack-bot/model"
)

// RenderProcesser process render
type RenderProcesser struct {
}

// Run render stage
func (obj *RenderProcesser) Run(stage *model.Stage, variables *map[string]interface{}) (string, error) {
	tmpl := template.New("temp")
	tmpl.Parse(stage.Template)
	var (
		data bytes.Buffer
	)
	err := tmpl.Execute(&data, *variables)
	if err != nil {
		return "", nil
	}
	return data.String(), nil
}
