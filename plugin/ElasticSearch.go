package plugin

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// ElasticSearchConfig plugin paramter
type ElasticSearchConfig struct {
	Variables map[string]string
	Query     string
}

// ElasticSearch plugin
func ElasticSearch(paramter []byte, variable *map[string]interface{}, output func(data interface{}), logger *zap.Logger) {
	var paramterConfig ElasticSearchConfig
	yaml.Unmarshal(paramter, &paramterConfig)
}
