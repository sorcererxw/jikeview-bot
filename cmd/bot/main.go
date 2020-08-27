package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/getsentry/sentry-go"
	tb "gopkg.in/tucnak/telebot.v2"

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
	if conf.IsAWSLambda {
		bot, err := tb.NewBot(tb.Settings{
			Token:       conf.BotToken,
			Synchronous: true,
			Reporter: func(err error) {
				if err.Error() == tb.ErrCouldNotUpdate.Error() {
					return
				}
				log.Print(err)
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		registerHandler(bot)
		lambda.Start(func(req events.APIGatewayProxyRequest) (err error) {
			var u tb.Update
			if err = json.Unmarshal([]byte(req.Body), &u); err == nil {
				bot.ProcessUpdate(u)
			}
			return
		})
	} else {
		var poller tb.Poller
		var removeWebhook bool
		if conf.AppEnv == "production" {
			removeWebhook = false
			poller = &tb.Webhook{
				Listen:   conf.Port,
				Endpoint: &tb.WebhookEndpoint{PublicURL: conf.WebHookEndpoint},
			}
		} else {
			removeWebhook = true
			poller = &tb.LongPoller{}
		}
		bot, err := tb.NewBot(tb.Settings{
			Token:  conf.BotToken,
			Poller: poller,
			Reporter: func(err error) {
				if err.Error() == tb.ErrCouldNotUpdate.Error() {
					return
				} else {
					log.Print(err)
				}
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		if removeWebhook {
			err := bot.RemoveWebhook()
			if err != nil {
				log.Fatal(err)
			}
		}
		registerHandler(bot)
		bot.Start()
	}
}
