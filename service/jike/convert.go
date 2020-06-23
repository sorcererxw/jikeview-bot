package jike

import (
	"fmt"
	"os"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/sorcererxw/jikeview-bot/util"
	"github.com/sorcererxw/jikeview-bot/util/log"
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
		limit := MaxCharacterLength - 100
		if p.hasMedia() {
			limit = MaxCaptionLength - 100
		}
		if len(content) > limit {
			content = content[:limit] + "......"
		}
		text += "\n" + content
	}
	text += fmt.Sprintf("\n\n<a href='%s'>查看全文</a>", p.GetUrl().GenerateMessageUrl())
	// video mode
	if p.GetVideo() != nil {
		video := p.GetVideo()
		mediaMeta, err := GetMediaMeta(p.GetUrl())
		if err != nil {
			return nil, err
		}
		filepath, err := util.DownloadAndFormatVideo(mediaMeta.Url)
		if err != nil {
			return nil, err
		}
		info, err := util.GetVideoInfo(filepath)
		if err != nil {
			return nil, err
		}
		file, err := os.Open(filepath)
		if err != nil {
			return nil, err
		}
		duration, err := strconv.ParseFloat(info.Format.Duration, 64)
		if err != nil {
			return nil, err
		}
		return &tb.Video{
			File:      tb.File{FileReader: file},
			Width:     info.Streams[0].Width,
			Height:    info.Streams[0].Height,
			Duration:  int(duration),
			Caption:   text,
			Thumbnail: &tb.Photo{File: tb.File{FileURL: video.ThumbnailUrl}},
		}, nil
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
		return &tb.Photo{
			File:      tb.File{FileURL: pic.PicUrl},
			Height:    pic.Height,
			Width:     pic.Width,
			Caption:   text,
			ParseMode: tb.ModeHTML,
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
