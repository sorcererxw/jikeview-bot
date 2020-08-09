package jike

import (
	"regexp"
)

type (
	Url struct {
		Type PostType
		ID   string
	}
)

func ParseUrl(url string) *Url {
	processors := []struct {
		matcher  *regexp.Regexp
		postType PostType
		idIndex  int
	}{
		{
			regexp.MustCompile("^(https?://)?web\\.(jellow\\.club|okjike\\.com)/post-detail/([0-9a-z]+)/originalPost(\\?.*)?$"),
			OriginalPost, 3,
		},
		{
			regexp.MustCompile("^(https?://)?web\\.(jellow\\.club|okjike\\.com)/message-detail/([0-9a-z]+)/officialMessage(\\?.*)?$"),
			OfficialMessage, 3,
		},
		{
			regexp.MustCompile("^(https?://)?(m|web)\\.(jellow\\.club|okjike\\.com)/originalPosts?/([0-9a-z]+)(\\?.*)?$"),
			OriginalPost, 4,
		},
		{
			regexp.MustCompile("^(https?://)?(m|web)\\.(jellow\\.club|okjike\\.com)/reposts?/([0-9a-z]+)(\\?.*)?$"),
			Repost, 4,
		},
		{
			regexp.MustCompile("^(https?://)?(m|web)\\.(jellow\\.club|okjike\\.com)/officialMessages?/([0-9a-z]+)(\\?.*)?$"),
			OfficialMessage, 4,
		},
	}
	for _, p := range processors {
		if p.matcher.MatchString(url) {
			return &Url{
				Type: p.postType,
				ID:   p.matcher.FindStringSubmatch(url)[p.idIndex],
			}
		}
	}
	return nil
}

func (u *Url) GenerateMessageUrl() string {
	if u.Type == OfficialMessage {
		return "https://m.okjike.com/officialMessages/" + u.ID
	}
	if u.Type == OriginalPost {
		return "https://m.okjike.com/originalPosts/" + u.ID
	}
	if u.Type == Repost {
		return "https://m.okjike.com/reposts/" + u.ID
	}
	return ""
}
