package conf

import (
	"log"
	"os"
)

var (
	BotToken        = os.Getenv("BOT_TOKEN")
	SentryDSN       = os.Getenv("SENTRY_DSN")
	AppEnv          = os.Getenv("APP_ENV")
	WebHookEndpoint = os.Getenv("WEB_HOOK_ENDPOINT")
	Port            = os.Getenv("PORT")
)

func init() {
	log.Printf("BotToken: %s", BotToken)
	log.Printf("SentryDSN: %s", SentryDSN)
	log.Printf("AppEnv: %s", AppEnv)
	log.Printf("WebHookEndpoint: %s", WebHookEndpoint)
	log.Printf("Port: %s", Port)
}
