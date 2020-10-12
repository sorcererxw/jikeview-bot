package jike

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/sorcererxw/jikeview-bot/util"
)

const MaxCharacterLength = 4096 // max character count in one telegram message
const MaxCaptionLength = 1024

func (p Post) hasMedia() bool {
	return p.GetVideo() != nil || p.GetAudio() != nil || len(p.Data.Pictures) > 0
}

// ConvertToTelegramPost convert jike post to telebot sendable
func (p *Post) ConvertToTelegramPost() (interface{}, error) {
	text := ""
	if p.Data.Topic != nil {
		text += "#" + util.RemoveAllSpaceAndPunctuation(p.Data.Topic.Content) + "\n"
	}

	if p.Data.Content != "" {
		content := p.Data.Content
		if len(content) > 0 {
			content = fmt.Sprintf("<b>@%s:</b> %s", p.Data.User.ScreenName, p.Data.Content)
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
	text += fmt.Sprintf("\n<a href='%s'>查看全文</a>", p.GetUrl().GenerateMessageUrl())
	// video mode
	if p.GetVideo() != nil {
		video := p.GetVideo()
		mediaMeta, err := GetMediaMeta(p.GetUrl())
		if err != nil {
			return nil, errors.WithStack(err)
		}
		filepath, err := util.DownloadAndFormatVideo(mediaMeta.Url)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		info, err := util.GetVideoInfo(filepath)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		file, err := os.Open(filepath)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		duration, err := strconv.ParseFloat(info.Format.Duration, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		result := tb.Video{
			File:     tb.File{FileReader: file},
			Width:    info.Streams[0].Width,
			Height:   info.Streams[0].Height,
			Duration: int(duration),
			Caption:  text,
		}
		if video.ThumbnailUrl != "" {
			thumbnail, err := util.DownloadImage(video.ThumbnailUrl)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			result.Thumbnail = &tb.Photo{File: tb.File{FileLocal: thumbnail}}
		}
		return &result, nil
	}
	// audio mode
	if p.GetAudio() != nil {
		audio := p.GetAudio()
		mediaMeta, err := GetMediaMeta(p.GetUrl())
		if err != nil {
			return nil, err
		}
		filePath, err := util.DownloadAndFormatAudio(mediaMeta.Url)
		if err != nil {
			return nil, err
		}
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		return &tb.Audio{
			File:      tb.File{FileReader: file},
			Duration:  audio.Duration,
			Caption:   text,
			Thumbnail: &tb.Photo{File: tb.File{FileURL: audio.OriginCoverUrl}},
			Title:     audio.Title,
			Performer: audio.Author,
		}, nil
	}
	// gif mode
	if len(p.Data.Pictures) > 0 && p.Data.Pictures[0].Format == "gif" {
		pic := p.Data.Pictures[0]
		return &tb.Animation{
			File:    tb.File{FileURL: pic.PicUrl},
			Height:  pic.Height,
			Width:   pic.Width,
			Caption: text,
		}, nil
	}
	// gallery mode
	if len(p.Data.Pictures) > 0 {
		var gallery []tb.InputMedia
		log.Println(len(p.Data.Pictures))
		log.Println(len(gallery))
		for _, pic := range p.Data.Pictures {
			log.Println(pic.PicUrl)
			if pic.Format == "gif" {
				continue
			}
			caption := ""
			if len(gallery) == 0 {
				caption = text
			}
			photo := &tb.Photo{
				File:      tb.File{FileURL: pic.PicUrl},
				Height:    pic.Height,
				Width:     pic.Width,
				Caption:   caption,
				ParseMode: tb.ModeHTML,
			}
			gallery = append(gallery, tb.InputMedia(photo))
		}
		return gallery, nil
	}
	// text mode
	return text, nil
}

func TryToConvertTelegramPost(url string) (interface{}, error) {
	jikeUrl := ParseUrl(url)
	if jikeUrl == nil {
		return nil, nil
	}
	jikePost, err := GetPost(jikeUrl)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sendable, err := jikePost.ConvertToTelegramPost()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return sendable, nil
}
