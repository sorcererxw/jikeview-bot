package jike

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	baseUrl = "https://api.ruguoapp.com/1.0"
)

var (
	httpClient = &http.Client{
		Timeout:   30 * time.Second,
		Transport: &http.Transport{MaxIdleConnsPerHost: 10},
	}
)

type PostType string

const (
	TypeOriginalPost    PostType = "ORIGINAL_POST"
	TypeOfficialMessage          = "OFFICIAL_MESSAGE"
	TypeRepost                   = "REPOST"
)

type (
	UrlContent struct {
		OriginalUrl string `json:"originalUrl"` // http://t.cn/EtoRw3B
		Title       string `json:"title"`       // t.cn
		Url         string `json:"url"`         // https://redirect.jike.ruguoapp.com?redirect=http%3A%2F%2Ft.cn%2FEtoRw3B
	}

	Picture struct {
		CropperPosX     float64 `json:"cropperPosX"`     // 0.5
		CropperPosY     float64 `json:"cropperPosY"`     // 0.5
		Format          string  `json:"format"`          // jpeg
		Height          int     `json:"height"`          // 461
		Width           int     `json:"width"`           // 440
		MiddlePicUrl    string  `json:"middlePicUrl"`    // https://cdn.ruguoapp.com/FlB30KmHjL_vVPjfXnJ_AAMVpK_z.jpg
		PicUrl          string  `json:"picUrl"`          // https://cdn.ruguoapp.com/FlB30KmHjL_vVPjfXnJ_AAMVpK_z.jpg
		SmallPicUrl     string  `json:"smallPicUrl"`     // https://cdn.ruguoapp.com/FlB30KmHjL_vVPjfXnJ_AAMVpK_z.jpg
		ThumbnailUrl    string  `json:"thumbnailUrl"`    // https://cdn.ruguoapp.com/FlB30KmHjL_vVPjfXnJ_AAMVpK_z.jpg
		WatermarkPicUrl string  `json:"watermarkPicUrl"` // https://cdn.ruguoapp.com/FlB30KmHjL_vVPjfXnJ_AAMVpK_z.jpg
	}

	PostTopic struct {
		Content string `json:"content"`
	}

	Video struct {
		Duration     int           `json:"duration"`
		Source       []interface{} `json:"source"`
		ThumbnailUrl string        `json:"thumbnailUrl"`
		Image        *Picture      `json:"image"`
		Type         string        `json:"type"`
	}

	Audio struct {
		ID             string `json:"id"`
		Type           string `json:"type"`
		Url            string `json:"url"`
		Author         string `json:"author"`
		CoverUrl       string `json:"coverUrl"`
		OriginCoverUrl string `json:"originCoverUrl"`
		Title          string `json:"title"`
		Duration       int    `json:"duration"`
	}

	LinkInfo struct {
		LinkUrl    string `json:"linkUrl"`    // http://video.weibo.com/show?fid=1034:4331341232195357
		PictureUrl string `json:"pictureUrl"` // https://pic-txcdn.ruguoapp.com/Fr71UsUFLBR8f_DQXYaeS2G5Se3d
		Source     string `json:"source"`     // 查看链接
		Title      string `json:"title"`      // 2018年最令人震撼的魔术表演
		Video      *Video `json:"video"`
	}

	UserInfo struct {
		ID         string `json:"id"`
		ScreenName string `json:"screenName"`
	}

	Post struct {
		Success *bool  `json:"success"`
		Error   string `json:"error"`
		Data    struct {
			ID          string        `json:"id"`
			Type        PostType      `json:"type"`
			Content     string        `json:"content"`
			UrlsInText  []*UrlContent `json:"urlsInText"`
			Status      string        `json:"status"`
			Pictures    []*Picture    `json:"pictures"`
			PictureUrls []*Picture    `json:"pictureUrls"`
			CreatedAt   string        `json:"createdAt"`
			MessageId   string        `json:"messageId"`
			Topic       *PostTopic    `json:"topic"`
			LinkUrl     string        `json:"linkUrl"`
			Video       *Video        `json:"video"`
			Audio       *Audio        `json:"audio"`
			LinkInfo    *LinkInfo     `json:"linkInfo"`
			User        *UserInfo     `json:"user"`
			Target      *Post         `json:"target"`
		} `json:"data"`
	}

	MediaMeta struct {
		MediaLink string `json:"mediaLink"` // https://www.instagram.com/p/BtYclm5govz/?isVideo=true
		Url       string `json:"url"`       // https://media-qncdn.ruguoapp.com/295e5f686-41b12.m3u8
	}
)

func GetPost(url *Url) (*Post, error) {
	tp := ""
	switch url.Type {
	case TypeOfficialMessage:
		tp = "officialMessages"
	case TypeOriginalPost:
		tp = "originalPosts"
	case TypeRepost:
		tp = "reposts"
	}
	log.Printf("GetPost %+v", url)
	req, err := http.NewRequest("GET", baseUrl+"/"+tp+"/get", nil)
	if err != nil {
		return nil, err
	}
	qs := req.URL.Query()
	qs.Add("id", url.ID)
	req.URL.RawQuery = qs.Encode()
	log.Print(req.URL.String())
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Print(err)
		}
	}()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp Post
	if err = json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, err
	}
	if resp.Success != nil && *resp.Success == false {
		return nil, fmt.Errorf(resp.Error)
	}
	return &resp, nil
}

func GetMediaMeta(url *Url) (*MediaMeta, error) {
	req, err := http.NewRequest("GET", baseUrl+"/mediaMeta/play", nil)
	if err != nil {
		return nil, err
	}
	qs := req.URL.Query()
	qs.Add("id", url.ID)
	qs.Add("type", string(url.Type))
	req.URL.RawQuery = qs.Encode()
	log.Print(req.URL.String())
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Print(err)
		}
	}()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp MediaMeta
	if err = json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *Post) GetVideo() *Video {
	if p.Data.Video != nil {
		return p.Data.Video
	}
	if p.Data.LinkInfo != nil {
		if p.Data.LinkInfo.Video != nil {
			return p.Data.LinkInfo.Video
		}
	}
	return nil
}

func (p *Post) GetUrl() *Url {
	return &Url{
		Type: p.Data.Type,
		ID:   p.Data.ID,
	}
}

func (p *Post) GetAudio() *Audio {
	return p.Data.Audio
}
