package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/getsentry/sentry-go"
	tb "gopkg.in/tucnak/telebot.v2"
	"mvdan.cc/xurls/v2"

	"github.com/sorcererxw/jikeview-bot/conf"
	"github.com/sorcererxw/jikeview-bot/service/cosmos"
	"github.com/sorcererxw/jikeview-bot/service/jike"
	"github.com/sorcererxw/jikeview-bot/service/jstore"
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

func ConvertToSendable(url string) (interface{}, error) {
	converters := []func(u string) (interface{}, error){
		jike.TryToConvertTelegramPost,
		jstore.TryToConvertToTelegramPost,
		cosmos.TryToConvertToTelegramPost,
	}
	for _, cvt := range converters {
		sendable, err := cvt(url)
		if err != nil {
			return nil, err
		}
		if sendable != nil {
			return sendable, nil
		}
	}
	return nil, nil
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
		b.Send(m.Sender, "hi\n\n使用帮助: /help")
	})

	b.Handle("/help", func(m *tb.Message) {
		SendSendable(b, m, strings.Join([]string{
			"将<b>即刻产品矩阵</b>内的外链发送给我，我就能将其解析成 Telegram 消息回复给您。",
			"",
			"目前支持:",
			"1. 即刻动态链接",
			"2. 即士多商品链接",
			"3. 小宇宙播客单集链接",
		}, "\n"))
	})

	b.Handle("/me", func(m *tb.Message) {
		log.Print(m.Chat.ID)
		log.Print(m.Sender.ID)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		urls := xurls.Strict().FindAllString(m.Text, -1)
		for _, url := range urls {
			sendable, err := ConvertToSendable(url)
			if err != nil {
				log.Println(err)
				continue
			}
			if sendable == nil {
				continue
			}
			err = SendSendable(b, m, sendable)
			switch err {
			case tb.ErrTooLarge:
				b.Send(m.Sender, fmt.Sprintf("%s 内文件过大，无法通过 Telegram 发送", url))
			default:
				b.Send(m.Sender, fmt.Sprintf("发送失败: %s", err.Error()))
				log.Error(err)
			}
		}
	})
	b.Start()
}
