package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/j75689/slack-bot/tool"

	elasticsearch6 "github.com/elastic/go-elasticsearch/v6"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// ElasticSearch6Config plugin paramter
type ElasticSearch6Config struct {
	Addresses []string
	User      string
	Pass      string
	Index     string
	Query     string
}

// ElasticSearch6 plugin
func ElasticSearch6(paramter []byte, variable *map[string]interface{}, output func(data interface{}), logger *zap.Logger) {
	logger.Info("Execute plugin [ ElasticSearch6 ]")
	var paramterConfig ElasticSearch6Config
	yaml.Unmarshal(paramter, &paramterConfig)

	cfg := elasticsearch6.Config{
		Addresses: paramterConfig.Addresses,
	}
	if paramterConfig.User != "" {
		cfg.Username = paramterConfig.User
	}
	if paramterConfig.Pass != "" {
		cfg.Password = paramterConfig.Pass
	}

	es, err := elasticsearch6.NewClient(cfg)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	query, _ := tool.ExcuteGoTemplate(paramterConfig.Query, *variable)
	var buf bytes.Buffer
	buf.WriteString(query)

	logger.Debug(fmt.Sprintf("Elastic Search query:\n%v", buf.String()))

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(paramterConfig.Index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer res.Body.Close()

	logger.Debug(fmt.Sprintf("Elastic Search result:\n%v", res.String()))

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		logger.Error(err.Error())
		return
	}
	output(result)
}
