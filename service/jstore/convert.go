package jstore

import (
	"fmt"

	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/sorcererxw/jikeview-bot/util/log"
)

const MaxCharacterLength = 4096 // max character count in one telegram message
const MaxCaptionLength = 1024

func (p *Product) hasMedia() bool {
	return len(p.Pictures) > 0
}

// ConvertToTelegramPost converts jstore product to telegram Sendable
func (p *Product) ConvertToTelegramPost() (interface{}, error) {
	text := ""
	text += fmt.Sprintf("<b>%s</b>\n", p.Title)

	if p.Review != "" {
		content := p.Review
		if len(content) > 0 {
			content = fmt.Sprintf("<b>@%s:</b> %s", p.User.Nickname, p.Review)
		}
		limit := MaxCharacterLength - 100
		if p.hasMedia() {
			limit = MaxCaptionLength - 100
		}
		if len([]rune(content)) > limit {
			content = string([]rune(content)[:limit]) + "......"
		}
		text += "\n" + content
		if len(content) > 0 {
			text = text + "\n"
		}
	}

	text += fmt.Sprintf("\n<a href='%s'>查看全文</a>", p.generateUrl())

	if len(p.Pictures) > 0 {
		var gallery []tb.InputMedia
		for _, pic := range p.Pictures {
			log.Println(pic.PicURL)
			if pic.Format == "gif" {
				continue
			}
			caption := ""
			if len(gallery) == 0 {
				caption = text
			}
			photo := &tb.Photo{
				File:      tb.File{FileURL: pic.PicURL},
				Height:    pic.Height,
				Width:     pic.Width,
				Caption:   caption,
				ParseMode: tb.ModeHTML,
			}
			gallery = append(gallery, tb.InputMedia(photo))
		}
		return gallery, nil
	}

	return text, nil
}

func TryToConvertToTelegramPost(url string) (interface{}, error) {
	u := ParseUrl(url)
	if u == nil {
		return nil, nil
	}
	p, err := GetProduct(u.ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetProduct")
	}
	sendable, err := p.ConvertToTelegramPost()
	if err != nil {
		return nil, errors.Wrap(err, "ConvertToTelegramPost")
	}
	return sendable, nil
}
