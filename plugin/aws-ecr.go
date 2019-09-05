package plugin

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// AWSEcrConfig plugin paramter
type AWSEcrConfig struct {
	Addresses []string
	User      string
	Pass      string
	Index     string
	Query     string
}

// AWSEcr plugin
func AWSEcr(paramter []byte, variable *map[string]interface{}, output func(data interface{}), logger *zap.Logger) {
	logger.Info("Execute plugin [ AWS Ecr ]")
	var paramterConfig AWSEcrConfig
	yaml.Unmarshal(paramter, &paramterConfig)

}
