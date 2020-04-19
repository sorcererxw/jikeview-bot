package conf

import "os"

var (
	BotToken        = os.Getenv("BOT_TOKEN")
	SentryDSN       = os.Getenv("SENTRY_DSN")
	AppEnv          = os.Getenv("APP_ENV")
	WebHookEndpoint = os.Getenv("WEB_HOOK_ENDPOINT")
	WebHookPort     = os.Getenv("WEB_HOOK_PORT")
)
