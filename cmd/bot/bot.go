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

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token: conf.BotToken,
		//URL:    "http://195.129.111.17:8012",
		Poller: &tb.LongPoller{},
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
		b.Reply(m, "hi")
	})

	b.Handle("/help", func(m *tb.Message) {
		b.Reply(m, "may I help U?")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		fmt.Println(m.Text)
		log.Printf("receive msg: %s", m.Text)
		urls := xurls.Strict().FindAllString(m.Text, -1)
		for _, url := range urls {
			log.Println(url)
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
