package main

import (
	"log"
	"net/http"

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
	if conf.AppEnv == "production" {
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
			return ctx.String(http.StatusOK, "healthy")
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
