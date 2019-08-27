package model

// HandlerConfig definition resource
type HandlerConfig struct {
	ID        int                    `yaml:"-,omitempty"`
	ConfigID  string                 `yaml:"configID"`
	Product   string                 `yaml:"product"`
	Variables map[string]interface{} `yaml:"variables"`
	Command   []string               `yaml:"command"`
	Stage     []*HandlerStageConfig  `yaml:"stage"`
}

// HandlerStageConfig definition stage resource
type HandlerStageConfig struct {
	Type     string                 `yaml:"type"` // render, action
	Plugin   string                 `yaml:"plugin"`
	Paramter map[string]interface{} `yaml:"paramter"`
	Output   map[string]string      `yaml:"output"`
	Template string                 `yaml:"template"` // for render
}
