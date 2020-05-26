package jike

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sorcererxw/jikeview-bot/util"
)

func TestGetPost(t *testing.T) {
	url := "https://m.jellow.club/originalPosts/5e9aa689ae00f00018fc586e?username=86cdd8bd-b8fc-472d-9240-f28358749211"
	parsedUrl := ParseUrl(url)
	post, err := GetPost(parsedUrl)
	if err != nil {
		t.Fatal(err)
	}
	if post.Data.ID != parsedUrl.ID {
		t.Fail()
	}
	if post.Data.Type != RawOriginalPost {
		t.Log(post.Data.Type)
		t.Fail()
	}
}

func TestGetDeletedPost(t *testing.T) {
	url := "https://m.jellow.club/originalPosts/5e9a9247266e310018cb3251"
	parsedUrl := ParseUrl(url)
	post, _ := GetPost(parsedUrl)
	assert.Nil(t, post)
}

func TestDownloadAndFormatVideo(t *testing.T) {
	url := "https://m.jellow.club/originalPosts/5e9aa689ae00f00018fc586e?username=86cdd8bd-b8fc-472d-9240-f28358749211"
	parsedUrl := ParseUrl(url)
	post, err := GetPost(parsedUrl)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(post.GetVideo().Type)
	meta, err := GetMediaMeta(post.GetUrl())
	t.Log(meta)
	if err != nil {
		t.Fatal(err)
	}
	videoFile, err := util.DownloadAndFormatVideo(meta.Url)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(videoFile)
}
