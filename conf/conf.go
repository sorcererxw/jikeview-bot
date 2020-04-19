package conf

import "os"

var (
	BotToken  = os.Getenv("BOT_TOKEN")
	SentryDSN = os.Getenv("SENTRY_DSN")
)
