package appruntime

// env variables
type env struct {
	AppName  string `envconfig:"APP_NAME" default:"slackbot"`
	Port     int16  `envconfig:"PORT" default:"6000"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"INFO"`
	LogPath  string `envconfig:"LOG_PATH" default:"./logs"`
	// Slack
	SlackVerificationToken string `envconfig:"SLACK_VERIFICATION_TOKEN"`
	SlackBotOauthToken     string `envconfig:"SLACK_BOT_OAUTH_TOKEN"`
}
