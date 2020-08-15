package main

import (
	"log"

	"github.com/getsentry/sentry-go"

	"github.com/sorcererxw/jikeview-bot/bot"
	"github.com/sorcererxw/jikeview-bot/conf"
)

func init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: conf.SentryDSN,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	bot.Bot.Start()
}
