package jstore

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	Picture struct {
		PicURL       string `json:"picUrl"`
		LargePicURL  string `json:"largePicUrl"`
		MiddlePicURL string `json:"middlePicUrl"`
		ThumbnailURL string `json:"thumbnailUrl"`
		Format       string `json:"format"`
		Width        int    `json:"width"`
		Height       int    `json:"height"`
	}

	Store struct {
		ID        string   `json:"id"`
		Name      string   `json:"name"`
		SelfNames []string `json:"selfNames"`
		UserID    string   `json:"userId"`
	}

	User struct {
		ID             string   `json:"id"`
		Nickname       string   `json:"nickname"`
		Store          *Store   `json:"store"`
		Avatar         *Picture `json:"avatar"`
		AvatarURL      string   `json:"avatarUrl"`
		FollowingCount int      `json:"followingCount"`
		FollowedCount  int      `json:"followedCount"`
	}

	Product struct {
		ID            string     `json:"id"`
		Type          string     `json:"type"`
		Title         string     `json:"title"`
		Review        string     `json:"review"`
		Shelf         string     `json:"shelf"`
		Picture       *Picture   `json:"picture"`
		Pictures      []*Picture `json:"pictures"`
		Store         *Store     `json:"store"`
		User          *User      `json:"user"`
		Vendor        string     `json:"vendor"`
		VendorName    string     `json:"vendorName"`
		VendorItemID  string     `json:"vendorItemId"`
		VendorItemURL string     `json:"vendorItemUrl"`
		VendorIcon    string     `json:"vendorIcon"`
		LikeCount     int        `json:"likeCount"`
		ViewCount     int        `json:"viewCount"`
		SharingURL    string     `json:"sharingUrl"`
		ShareCount    int        `json:"shareCount"`
		CommentCount  int        `json:"commentCount"`
		CollectCount  int        `json:"collectCount"`
	}
)

const (
	baseURL = "https://api.jstore.site"
)

var (
	httpClient = &http.Client{
		Timeout:   30 * time.Second,
		Transport: &http.Transport{MaxIdleConnsPerHost: 10},
	}
)

func GetProduct(ID string) (*Product, error) {
	req, err := http.NewRequest("GET", baseURL+"/v1/products/get", nil)
	if err != nil {
		return nil, err
	}
	qs := req.URL.Query()
	qs.Add("id", ID)
	req.URL.RawQuery = qs.Encode()
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp Product
	if err = json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
