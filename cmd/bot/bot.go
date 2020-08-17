package main

import (
	"fmt"
	"log"
	"reflect"

	tb "gopkg.in/tucnak/telebot.v2"
	"mvdan.cc/xurls/v2"

	"github.com/sorcererxw/jikeview-bot/service/jike"
	"github.com/sorcererxw/jikeview-bot/service/jstore"
)

func SendSendable(bot *tb.Bot, m *tb.Message, sendable interface{}) error {
	rt := reflect.TypeOf(sendable)
	switch rt.Kind() {
	case reflect.Array, reflect.Slice:
		_, err := bot.SendAlbum(m.Sender, sendable.([]tb.InputMedia), &tb.SendOptions{
			ParseMode: tb.ModeHTML,
		})
		return err
	default:
		_, err := bot.Send(m.Sender, sendable, &tb.SendOptions{
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

func registerHandler(bot *tb.Bot) {
	bot.Handle("/start", func(m *tb.Message) {
		if err := SendSendable(bot, m, "hi\n\n使用帮助: /help"); err != nil {
			log.Print(err)
		}
	})

	bot.Handle("/help", func(m *tb.Message) {
		if err := SendSendable(bot, m, "将<b>即刻/即士多</b>内的消息链接发送给我，我就能将其解析成 Telegram 消息回复给您。"); err != nil {
			log.Print(err)
		}
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		urls := xurls.Strict().FindAllString(m.Text, -1)
		if len(urls) == 0 {
			bot.Send(m.Sender, "未识别到有效链接", &tb.SendOptions{
				ReplyTo: m,
			})
			return
		}
		for _, url := range urls {
			sendable, err := ConvertToSendable(url)
			if err != nil {
				log.Println(err)
				bot.Send(m.Sender, fmt.Sprintf("无法处理链接: %s", url), &tb.SendOptions{
					ReplyTo: m,
				})
				continue
			}
			if sendable == nil {
				bot.Send(m.Sender, fmt.Sprintf("无法转换链接: %s", url), &tb.SendOptions{
					ReplyTo: m,
				})
				continue
			}
			if err := SendSendable(bot, m, sendable); err != nil {
				switch err {
				case tb.ErrTooLarge:
					bot.Send(m.Sender, fmt.Sprintf("%s 内文件过大，无法通过 Telegram 发送", url), &tb.SendOptions{
						ReplyTo: m,
					})
				default:
					bot.Send(m.Sender, fmt.Sprintf("发送失败: %s", err.Error()), &tb.SendOptions{
						ReplyTo: m,
					})
					log.Print(err)
				}
			}
		}
	})
}
