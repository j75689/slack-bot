package plugin

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/j75689/slack-bot/appruntime"
	"go.uber.org/zap"
)

// Plugin ...
type Plugin func(paramter []byte, variable *map[string]interface{}, output func(data interface{}), logger *zap.Logger)

// Pool ...
type Pool struct {
	*sync.Map
}

// Execute ...
func (pool *Pool) Execute(plugin string, paramter []byte, variable *map[string]interface{}, output func(data interface{})) error {
	if p, ok := pool.Load(plugin); ok {
		if pluginfunc, ok := p.(Plugin); ok {
			pluginfunc(paramter, variable, output, appruntime.Logger)
			return nil
		}
		return fmt.Errorf("plugin [%s] load error", plugin)
	}
	return fmt.Errorf("plugin [%s] not found", plugin)
}

// Load ...
func Load() (pool *Pool) {
	pool = &Pool{
		Map: &sync.Map{},
	}
	prepare(
		pool,
		ElasticSearch6,
	)
	return
}

func prepare(pool *Pool, plugins ...Plugin) {
	for _, plugin := range plugins {
		pluginName := runtime.FuncForPC(reflect.ValueOf(plugin).Pointer()).Name()
		pluginName = pluginName[strings.LastIndex(pluginName, ".")+1:]
		pool.Store(pluginName, plugin)
	}
}
