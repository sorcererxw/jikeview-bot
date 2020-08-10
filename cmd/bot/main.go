package main

import (
	"github.com/getsentry/sentry-go"

	"github.com/sorcererxw/jikeview-bot/bot"
	"github.com/sorcererxw/jikeview-bot/conf"
	"github.com/sorcererxw/jikeview-bot/util/log"
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
