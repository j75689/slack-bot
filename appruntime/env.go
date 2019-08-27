package appruntime

// env variables
type env struct {
	// Basic Info
	AppName  string `envconfig:"APP_NAME" default:"slackbot"`
	Port     int16  `envconfig:"PORT" default:"6000"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"INFO"`
	LogPath  string `envconfig:"LOG_PATH" default:"./logs"`
	// Slack
	SlackVerificationToken string `envconfig:"SLACK_VERIFICATION_TOKEN"`
	SlackBotOauthToken     string `envconfig:"SLACK_BOT_OAUTH_TOKEN"`
	// DB
	DBDriver string `envconfig:"DB_DRIVER" default:"Bolt"`
	DBName   string `envconfig:"DB_NAME"`
	DBHost   string `envconfig:"DB_HOST" default:"slackbot.db"`
	DBPort   string `envconfig:"DB_PORT"`
	DBUser   string `envconfig:"DB_USER"`
	DBPass   string `envconfig:"DB_PASS"`
}
