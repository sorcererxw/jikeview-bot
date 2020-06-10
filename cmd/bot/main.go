package main

import (
	"fmt"
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

func SendSendable(b *tb.Bot, m *tb.Message, sendable interface{}) error {
	rt := reflect.TypeOf(sendable)
	switch rt.Kind() {
	case reflect.Array, reflect.Slice:
		_, err := b.SendAlbum(m.Sender, sendable.([]tb.InputMedia), &tb.SendOptions{
			ParseMode: tb.ModeHTML,
		})
		return err
	default:
		_, err := b.Send(m.Sender, sendable, &tb.SendOptions{
			ParseMode: tb.ModeHTML,
		})
		return err
	}
}

func main() {
	var poller tb.Poller
	var removeWebhook bool
	if conf.AppEnv == "production" {
		removeWebhook = false
		poller = &tb.Webhook{
			Listen:   conf.WebHookPort,
			Endpoint: &tb.WebhookEndpoint{PublicURL: conf.WebHookEndpoint},
		}
	} else {
		removeWebhook = true
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
	if removeWebhook {
		err := b.RemoveWebhook()
		if err != nil {
			log.Fatal(err)
		}
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
			err = SendSendable(b, m, sendable)
			switch err {
			case tb.ErrTooLarge:
				b.Send(m.Sender, fmt.Sprintf("%s 内文件过大，无法通过 Telegram 发送", url))
			}
		}
	})
	b.Start()
}
