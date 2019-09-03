package plugin

import (
	"fmt"
	"io/ioutil"
	"plugin"
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
		}
		return nil
	}
	return fmt.Errorf("plugin [%s] not found", plugin)
}

// Load ...
func Load(path string) (pool *Pool) {
	pool = &Pool{
		Map: &sync.Map{},
	}
	// fix Path
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// read files
	files, err := ioutil.ReadDir(path)
	if err != nil {
		appruntime.Logger.Error(fmt.Sprintf("load plugin error: " + err.Error()))
		return
	}

	// load
	for _, f := range files {
		if !f.IsDir() {
			var runFuncName = f.Name()

			appruntime.Logger.Info(fmt.Sprintf("loading plugin [%s]", runFuncName))

			if !strings.HasSuffix(f.Name(), ".so") {
				continue
			}

			if strings.LastIndexAny(runFuncName, ".") > -1 {
				runFuncName = runFuncName[0:strings.LastIndexAny(runFuncName, ".")]
			}

			p, err := plugin.Open(path + f.Name())
			if err != nil {
				appruntime.Logger.Error(fmt.Sprintf("open plugin [%s] error: %v", path+f.Name(), err))
				continue
			}

			function, err := p.Lookup(runFuncName)
			if err != nil {
				appruntime.Logger.Error(fmt.Sprintf("lookup func [%s] error: %v", runFuncName, err))
				continue
			}

			if f, ok := function.(Plugin); ok {
				pool.Store(runFuncName, f)
			} else {
				appruntime.Logger.Error(fmt.Sprintf("reflect func [%s] error: %v", runFuncName, err))
			}

		}
	}

	return
}
