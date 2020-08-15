package bot

import (
	"fmt"
	"log"
	"reflect"

	tb "gopkg.in/tucnak/telebot.v2"
	"mvdan.cc/xurls/v2"

	"github.com/sorcererxw/jikeview-bot/conf"
	"github.com/sorcererxw/jikeview-bot/service/jike"
	"github.com/sorcererxw/jikeview-bot/service/jstore"
)

var Bot *tb.Bot

func SendSendable(m *tb.Message, sendable interface{}) error {
	rt := reflect.TypeOf(sendable)
	switch rt.Kind() {
	case reflect.Array, reflect.Slice:
		_, err := Bot.SendAlbum(m.Sender, sendable.([]tb.InputMedia), &tb.SendOptions{
			ParseMode: tb.ModeHTML,
		})
		return err
	default:
		_, err := Bot.Send(m.Sender, sendable, &tb.SendOptions{
			ParseMode: tb.ModeHTML,
		})
		return err
	}
}

func ConvertToSendable(url string) (interface{}, error) {
	converters := []func(u string) (interface{}, error){
		jike.TryToConvertTelegramPost,
		jstore.TryToConvertToTelegramPost,
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

func init() {
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
	var err error
	Bot, err = tb.NewBot(tb.Settings{
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
		return
	}
	if removeWebhook {
		err := Bot.RemoveWebhook()
		if err != nil {
			log.Fatal(err)
		}
	}

	Bot.Handle("/start", func(m *tb.Message) {
		Bot.Send(m.Sender, "hi\n\n使用帮助: /help")
	})

	Bot.Handle("/help", func(m *tb.Message) {
		SendSendable(m, "将<Bot>即刻/即士多</Bot>内的消息链接发送给我，我就能将其解析成 Telegram 消息回复给您。")
	})

	Bot.Handle(tb.OnText, func(m *tb.Message) {
		urls := xurls.Strict().FindAllString(m.Text, -1)
		if len(urls) == 0 {
			Bot.Send(m.Sender, "未识别到有效链接", &tb.SendOptions{
				ReplyTo: m,
			})
			return
		}
		for _, url := range urls {
			sendable, err := ConvertToSendable(url)
			if err != nil {
				log.Println(err)
				continue
			}
			if sendable == nil {
				Bot.Send(m.Sender, fmt.Sprintf("无法转换链接: %s", url), &tb.SendOptions{
					ReplyTo: m,
				})
				continue
			}
			if err := SendSendable(m, sendable); err != nil {
				switch err {
				case tb.ErrTooLarge:
					Bot.Send(m.Sender, fmt.Sprintf("%s 内文件过大，无法通过 Telegram 发送", url), &tb.SendOptions{
						ReplyTo: m,
					})
				default:
					Bot.Send(m.Sender, fmt.Sprintf("发送失败: %s", err.Error()), &tb.SendOptions{
						ReplyTo: m,
					})
					log.Print(err)
				}
			}
		}
	})
}
