package jike

import (
	"regexp"
)

type (
	Url struct {
		Type SimplePostType
		ID   string
	}
)

func ParseUrl(url string) *Url {
	WebOriginalPost := regexp.MustCompile("^(https?://)?web\\.(jellow\\.club|okjike\\.com)/post-detail/([0-9a-z]+)/originalPost(\\?.*)?$")
	WebOfficialMessage := regexp.MustCompile("^(https?://)?web\\.(jellow\\.club|okjike\\.com)/message-detail/([0-9a-z]+)/officialMessage(\\?.*)?$")
	MobileOriginalPost := regexp.MustCompile("^(https?://)?m\\.(jellow\\.club|okjike\\.com)/originalPosts/([0-9a-z]+)(\\?.*)?$")
	MobileOfficialMessage := regexp.MustCompile("^(https?://)?m\\.(jellow\\.club|okjike\\.com)/officialMessages/([0-9a-z]+)(\\?.*)?$")
	switch {
	case WebOriginalPost.MatchString(url):
		return &Url{OriginalPost, WebOriginalPost.FindStringSubmatch(url)[3]}
	case WebOfficialMessage.MatchString(url):
		return &Url{OfficialMessage, WebOfficialMessage.FindStringSubmatch(url)[3]}
	case MobileOriginalPost.MatchString(url):
		return &Url{OriginalPost, MobileOriginalPost.FindStringSubmatch(url)[3]}
	case MobileOfficialMessage.MatchString(url):
		return &Url{OfficialMessage, MobileOfficialMessage.FindStringSubmatch(url)[3]}
	default:
		return nil
	}
}

func (u *Url) GenerateMessageUrl() string {
	if u.Type == OfficialMessage {
		return "https://m.okjike.com/officialMessages" + u.ID
	}
	if u.Type == OriginalPost {
		return "https://m.okjike.com/originalPosts/" + u.ID
	}
	return ""
}
