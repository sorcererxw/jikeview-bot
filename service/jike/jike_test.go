package jike

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sorcererxw/jikeview-bot/util"
)

func TestGetPost(t *testing.T) {
	t.Skip("skip TestGetPost")
	url := "https://api.ruguoapp.com/originalPosts/5e9aa689ae00f00018fc586e"
	parsedUrl := ParseUrl(url)
	post, err := GetPost(parsedUrl)
	assert.Nil(t, err)
	assert.Equal(t, parsedUrl.ID, post.Data.ID)
	assert.Equal(t, TypeOriginalPost, post.Data.Type)
}

func TestGetDeletedPost(t *testing.T) {
	t.Skip("skip TestGetDeletedPost")
	url := "https://api.ruguoapp.com/originalPosts/5e9a9247266e310018cb3251"
	parsedUrl := ParseUrl(url)
	post, _ := GetPost(parsedUrl)
	assert.Nil(t, post)
}

func TestDownloadAndFormatVideo(t *testing.T) {
	t.Skip("skip TestDownloadAndFormatVideo")
	url := "https://api.ruguoapp.com/originalPosts/5e9aa689ae00f00018fc586e?username=86cdd8bd-b8fc-472d-9240-f28358749211"
	parsedUrl := ParseUrl(url)
	post, err := GetPost(parsedUrl)
	assert.Nil(t, err)
	meta, err := GetMediaMeta(post.GetUrl())
	assert.Nil(t, err)
	videoFile, err := util.DownloadAndFormatVideo(meta.Url)
	assert.Nil(t, err)
	t.Log(videoFile)
}
