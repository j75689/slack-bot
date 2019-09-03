package model

// SlackBotConfig define basic fields
type SlackBotConfig struct {
	ID       int       `yaml:"-"`
	Version  string    `yaml:"version"`
	Kind     string    `yaml:"kind"` // Project , Message , ...
	MetaData *MetaData `yaml:"metadata"`
	Task     *Task     `yaml:"task,omitempty"`
}

// MetaData define resource name
type MetaData struct {
	Name    string            `yaml:"name"`
	Project string            `yaml:"project,omitempty"`
	Labels  map[string]string `yaml:"labels,omitempty"`
}

// Task define action
type Task struct {
	Variables map[string]interface{} `yaml:"variables"`
	Command   []string               `yaml:"command"`
	Stage     []*Stage               `yaml:"stage"`
}

// Stage define execute
type Stage struct {
	Type     string                 `yaml:"type"` // render, action
	Plugin   string                 `yaml:"plugin,omitempty"`
	Paramter map[string]interface{} `yaml:"paramter,omitempty"`
	Output   string                 `yaml:"output,omitempty"`   // output value name
	Template string                 `yaml:"template,omitempty"` // for render
}
