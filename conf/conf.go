package conf

import (
	"os"

	"github.com/sorcererxw/jikeview-bot/util/log"
)

var (
	BotToken        = os.Getenv("BOT_TOKEN")
	SentryDSN       = os.Getenv("SENTRY_DSN")
	AppEnv          = os.Getenv("APP_ENV")
	WebHookEndpoint = os.Getenv("WEB_HOOK_ENDPOINT")
	WebHookPort     = os.Getenv("WEB_HOOK_PORT")
)

func init() {
	log.Printf("BotToken: %s", BotToken)
	log.Printf("SentryDSN: %s", SentryDSN)
	log.Printf("AppEnv: %s", AppEnv)
	log.Printf("WebHookEndpoint: %s", WebHookEndpoint)
	log.Printf("WebHookPort: %s", WebHookPort)
}
