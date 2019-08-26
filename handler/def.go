package handler

// StageType ...
type StageType string

const (
	// stage type
	actionTag StageType = "action"
	renderTag StageType = "render"
)

// Handler handle slack message
type Handler interface {
	Do() (string, error)
	DryRun() (string, error)
}

// StageRunner run stage action
type StageRunner interface {
	Run(variables map[string]interface{}) (string, error)
}
