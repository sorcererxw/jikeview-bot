package conf

import (
	"log"
	"os"
	"regexp"
)

var (
	BotToken        = os.Getenv("BOT_TOKEN")
	SentryDSN       = os.Getenv("SENTRY_DSN")
	AppEnv          = os.Getenv("APP_ENV")
	WebHookEndpoint = os.Getenv("WEB_HOOK_ENDPOINT")
	Port            = os.Getenv("PORT")
	IsAWSLambda     = regexp.MustCompile("^AWS_Lambda_").MatchString(os.Getenv("AWS_EXECUTION_ENV"))
)

func init() {
	log.Printf("BotToken: %s", BotToken)
	log.Printf("SentryDSN: %s", SentryDSN)
	log.Printf("AppEnv: %s", AppEnv)
	log.Printf("WebHookEndpoint: %s", WebHookEndpoint)
	log.Printf("Port: %s", Port)
}
