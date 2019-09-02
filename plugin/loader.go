package plugin

import (
	"go.uber.org/zap"
)

// Plugin ...
type Plugin func(paramter []byte, variable *map[string]interface{}, logger zap.Logger)

// Load ...
func Load(path string) {

}
