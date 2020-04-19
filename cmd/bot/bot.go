package main

import (
	"reflect"

	"github.com/getsentry/sentry-go"
	tb "gopkg.in/tucnak/telebot.v2"
	"mvdan.cc/xurls/v2"

	"github.com/sorcererxw/jikeview-bot/conf"
	"github.com/sorcererxw/jikeview-bot/service/jike"
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
	var poller tb.Poller
	if conf.AppEnv == "production" {
		poller = &tb.Webhook{
			Listen:   conf.WebHookPort,
			Endpoint: &tb.WebhookEndpoint{PublicURL: conf.WebHookEndpoint},
		}
	} else {
		poller = &tb.LongPoller{}
	}
	b, err := tb.NewBot(tb.Settings{
		Token:  conf.BotToken,
		Poller: poller,
		Reporter: func(err error) {
			if err.Error() == tb.ErrCouldNotUpdate.Error() {
				return
			} else {
				log.Error(err)
			}
		},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "hi")
	})

	b.Handle("/help", func(m *tb.Message) {
		b.Send(m.Sender, "将即刻App内的消息链接发送给我，我就能将其解析成 Telegram 消息回复给您。")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		urls := xurls.Strict().FindAllString(m.Text, -1)
		for _, url := range urls {
			jikeUrl := jike.ParseUrl(url)
			if jikeUrl == nil {
				continue
			}
			jikePost, err := jike.GetPost(jikeUrl)
			if err != nil {
				log.Error(err)
				return
			}
			sendable, err := jikePost.ConvertToTelegramPost()
			if err != nil {
				log.Error(err)
				return
			}
			rt := reflect.TypeOf(sendable)
			switch rt.Kind() {
			case reflect.Array, reflect.Slice:
				_, err := b.SendAlbum(m.Sender, sendable.([]tb.InputMedia), &tb.SendOptions{
					ParseMode: tb.ModeHTML,
				})
				if err != nil {
					log.Error(err)
				}
			default:
				_, err := b.Send(m.Sender, sendable, &tb.SendOptions{
					ParseMode: tb.ModeHTML,
				})
				if err != nil {
					log.Error(err)
				}
			}
		}
	})

	b.Start()
}
