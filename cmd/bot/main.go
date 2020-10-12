package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
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
	} else if conf.AppEnv == "production" {
		bot, err := tb.NewBot(tb.Settings{
			Token: conf.BotToken,
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
		registerHandler(bot)
		e := echo.New()
		e.GET("/health", func(ctx echo.Context) error {
			return ctx.NoContent(http.StatusOK)
		})
		e.POST("/bot", func(ctx echo.Context) error {
			var u tb.Update
			err := ctx.Bind(&u)
			if err != nil {
				return err
			}
			bot.ProcessUpdate(u)
			return ctx.NoContent(http.StatusOK)
		})
		if err := e.Start(":" + conf.Port); err != nil {
			panic(err)
		}
	} else {
		bot, err := tb.NewBot(tb.Settings{
			Token:  conf.BotToken,
			Poller: &tb.LongPoller{},
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
		if err := bot.RemoveWebhook(); err != nil {
			log.Fatal(err)
		}
		registerHandler(bot)
		bot.Start()
	}
}
