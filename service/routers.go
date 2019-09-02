package service

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/j75689/slack-bot/model"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/j75689/slack-bot/appruntime"
	"github.com/j75689/slack-bot/manager"
	"github.com/j75689/slack-bot/service/middleware"
	"github.com/j75689/slack-bot/service/modules"
	"github.com/nlopes/slack"
)

// InitRouter initial routers
func InitRouter() (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	router.HandleMethodNotAllowed = true
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(ginzap.Ginzap(appruntime.Logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(appruntime.Logger, true))
	router.NoRoute(middleware.NoRouteHandler())
	router.NoMethod(middleware.NoMethodHandler())

	register(router)

	return
}

func register(app *gin.Engine) {
	// health check
	health := app.Group("/health")
	{
		health.Any("", modules.HandleHealthCheck())
	}
	management := manager.NewManagement()
	initConfig(management)
	// slack hook
	slackhook := app.Group("/slack")
	slackhook.Use(middleware.VerifyProjectMiddleware(management))
	{
		api := slack.New(appruntime.Env.SlackBotOauthToken)
		authResp, err := api.AuthTest()
		if err != nil {
			appruntime.Logger.Fatal(err.Error())
		}
		botID := authResp.UserID
		slackhook.POST("/:project/events-endpoint", modules.HandleSlackEvent(api, botID, management))
		slackhook.Any("/:project/interactive-endpoint", modules.HandleSlackInteractive(management))
	}
	// debug
	debug := app.Group("/debug")
	{
		debug.POST("/dryrun", modules.HandleDryRun(management))
	}
	// api
	api := app.Group("/api/v1")
	{
		api.POST("/config/apply", modules.HandleApplyConfig(management))
		api.POST("/config/delete", modules.HandleDeleteConfig(management))
	}
}

func initConfig(management *manager.Management) {
	appruntime.DB.FindAll(func(project, kind, key string, data []byte) {
		appruntime.Logger.Info(fmt.Sprintf("load config project:[%s] kind:[%s] name:[%s]", project, kind, key))
		var config model.SlackBotConfig
		appruntime.Logger.Debug("unmarshal config ... \n" + string(data))
		if err := yaml.Unmarshal(data, &config); err != nil {
			appruntime.Logger.Error("load config error: " + err.Error())
			return
		}

		appruntime.Logger.Debug("find config manager ...")
		if ok, manager := management.Get(config.Kind); ok {
			appruntime.Logger.Debug("register config ...")
			// register new cmd
			_, err := manager.Register(&config)
			if err != nil {
				appruntime.Logger.Error("register config error: " + err.Error())
			}
			appruntime.Logger.Debug("register config complete")
		}
	})
}
