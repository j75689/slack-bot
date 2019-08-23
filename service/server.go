package service

import (
	"fmt"

	"github.com/Invisibi-nd/slack-bot/appruntime"
)

// Start Server
func Start() {
	appruntime.Logger.Info(fmt.Sprintf("Server Start on localhost:%d", appruntime.Env.Port))
	appruntime.Logger.Error((InitRouter().Run(fmt.Sprintf(":%d", appruntime.Env.Port))).Error())
}

// Shutdown Server
func Shutdown() {
	// ...
}
