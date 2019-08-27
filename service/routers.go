package service

import (
	"time"

	"github.com/Invisibi-nd/slack-bot/appruntime"
	"github.com/Invisibi-nd/slack-bot/handler"
	"github.com/Invisibi-nd/slack-bot/service/middleware"
	"github.com/Invisibi-nd/slack-bot/service/modules"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
)

// InitRouter initial routers
func InitRouter() (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(ginzap.Ginzap(appruntime.Logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(appruntime.Logger, true))
	router.NoMethod(middleware.NoMethodHandler())
	router.NoRoute(middleware.NoRouteHandler())
	registerAPI(router)

	return
}

func registerAPI(app *gin.Engine) {
	// health check
	health := app.Group("/health")
	{
		health.Any("", modules.HandleHealthCheck())
	}
	manager := handler.NewManager()
	// slack hook
	slackhook := app.Group("/slack")
	{
		api := slack.New(appruntime.Env.SlackBotOauthToken)
		authResp, err := api.AuthTest()
		if err != nil {
			appruntime.Logger.Fatal(err.Error())
		}
		botID := authResp.UserID
		slackhook.POST("/events-endpoint", modules.HandleSlackEvent(api, botID))
		slackhook.Any("/interactive-endpoint", modules.HandleSlackInteractive())
	}
	// Debug
	debug := app.Group("/debug")
	{
		debug.POST("/dryrun", modules.HandleDryRun(manager))
	}
}
