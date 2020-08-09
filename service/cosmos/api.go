package cosmos

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Episode struct {
	Type        string      `json:"type"` // example: EPISODE
	Eid         string      `json:"eid"`
	Pid         string      `json:"pid"`
	Title       string      `json:"title"`
	Subtitle    interface{} `json:"subtitle"`
	Shownotes   string      `json:"shownotes"`
	Description string      `json:"description"`
	Image       struct {
		PicURL       string `json:"picUrl"`
		LargePicURL  string `json:"largePicUrl"`
		MiddlePicURL string `json:"middlePicUrl"`
		SmallPicURL  string `json:"smallPicUrl"`
		ThumbnailURL string `json:"thumbnailUrl"`
	} `json:"image"`
	Enclosure struct {
		URL string `json:"url"`
	} `json:"enclosure"`
	ClapCount    int       `json:"clapCount"`
	CommentCount int       `json:"commentCount"`
	PlayCount    int       `json:"playCount"`
	PubDate      time.Time `json:"pubDate"`
	Duration     int       `json:"duration"`
	Podcast      struct {
		Type              string `json:"type"`
		Pid               string `json:"pid"`
		Title             string `json:"title"`
		Author            string `json:"author"`
		Description       string `json:"description"`
		SubscriptionCount int    `json:"subscriptionCount"`
		Image             struct {
			PicURL       string `json:"picUrl"`
			LargePicURL  string `json:"largePicUrl"`
			MiddlePicURL string `json:"middlePicUrl"`
			SmallPicURL  string `json:"smallPicUrl"`
			ThumbnailURL string `json:"thumbnailUrl"`
		} `json:"image"`
		LatestEpisodePubDate time.Time `json:"latestEpisodePubDate"`
		SubscriptionStatus   string    `json:"subscriptionStatus"`
		SubscriptionPush     bool      `json:"subscriptionPush"`
		Status               string    `json:"status"`
		Permissions          []struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		} `json:"permissions"`
	} `json:"podcast"`
	Permissions []struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"permissions"`
	CommentPermission string `json:"commentPermission"`
}

const (
	baseURL = "https://www.xiaoyuzhoufm.com"
)

var (
	httpClient = &http.Client{
		Timeout:   30 * time.Second,
		Transport: &http.Transport{MaxIdleConnsPerHost: 10},
	}
)

func GetEpisode(ID string) (*Episode, error) {
	req, err := http.NewRequest("GET", baseURL+"/api/episodes/"+ID, nil)
	if err != nil {
		return nil, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp Episode
	if err = json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
