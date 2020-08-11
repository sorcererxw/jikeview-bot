package cosmos

import (
	"fmt"

	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"jaytaylor.com/html2text"

	"github.com/sorcererxw/jikeview-bot/util/log"
)

const MaxCaptionLength = 1024

// ConvertToTelegramPost converts jstore product to telegram Sendable
func (p *Episode) ConvertToTelegramPost() (interface{}, error) {
	var err error

	text := ""
	text += fmt.Sprintf("<b>%s</b>\n", p.Title)

	content := ""
	if p.Shownotes != "" {
		content = p.Shownotes
	} else if p.Description != "" {
		content = p.Description
	}
	content, err = html2text.FromString(content)
	if err != nil {
		return nil, err
	}
	log.Print(content)
	if content != "" {
		limit := MaxCaptionLength - 100
		if len([]rune(content)) > limit {
			content = string([]rune(content)[:limit]) + "......"
		}
		text += "\n" + content
		if len(content) > 0 {
			text = text + "\n"
		}
	}

	text += fmt.Sprintf("\n<a href='%s'>打开小宇宙</a>", p.generateUrl())

	log.Print(p.Enclosure.URL)

	//audioFile, err := util.DownloadAndFormatAudio(p.Enclosure.URL)
	//if err != nil {
	//	return nil,err
	//}
	return &tb.Audio{
		File: tb.File{
			FileURL: p.Enclosure.URL,
		},
		Caption: text,
		Thumbnail: &tb.Photo{
			File: tb.File{
				FileURL: p.Image.PicURL,
			},
		},
		Title:     p.Title,
		Performer: p.Podcast.Author,
		MIME:      "audio/mpeg",
	}, nil
}

func TryToConvertToTelegramPost(url string) (interface{}, error) {
	u := ParseUrl(url)
	if u == nil {
		return nil, nil
	}
	p, err := GetEpisode(u.ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetEpisode")
	}
	sendable, err := p.ConvertToTelegramPost()
	if err != nil {
		return nil, errors.Wrap(err, "ConvertToTelegramPost")
	}
	return sendable, nil
}
