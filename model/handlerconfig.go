package model

// HandlerConfig definition resource
type HandlerConfig struct {
	ID        int                    `yaml:"-,omitempty"`
	Version   string                 `yaml:"version"`
	ConfigID  string                 `yaml:"id"`
	Project   string                 `yaml:"project"`
	Variables map[string]interface{} `yaml:"variables"`
	Command   []string               `yaml:"command"`
	Stage     []*HandlerStageConfig  `yaml:"stage"`
}

// HandlerStageConfig definition stage resource
type HandlerStageConfig struct {
	Type     string                 `yaml:"type"` // render, action
	Plugin   string                 `yaml:"plugin,omitempty"`
	Paramter map[string]interface{} `yaml:"paramter,omitempty"`
	Output   map[string]string      `yaml:"output,omitempty"`
	Template string                 `yaml:"template,omitempty"` // for render
}
