package jike

import (
	"regexp"
)

type (
	URL struct {
		Type PostType
		ID   string
	}
)

func ParseURL(url string) *URL {
	processors := []struct {
		matcher  *regexp.Regexp
		postType PostType
		idIndex  int
	}{
		{
			regexp.MustCompile(`^(https?://)?web\.(jellow\.club|okjike\.com)/post-detail/([0-9a-z]+)/originalPost(\?.*)?$`),
			TypeOriginalPost, 3,
		},
		{
			regexp.MustCompile(`^(https?://)?web\.(jellow\.club|okjike\.com)/message-detail/([0-9a-z]+)/officialMessage(\?.*)?$`),
			TypeOfficialMessage, 3,
		},
		{
			regexp.MustCompile(`^(https?://)?(m|web)\.(jellow\.club|okjike\.com)/originalPosts?/([0-9a-z]+)(\?.*)?$`),
			TypeOriginalPost, 4,
		},
		{
			regexp.MustCompile(`^(https?://)?(m|web)\.(jellow\.club|okjike\.com)/reposts?/([0-9a-z]+)(\?.*)?$`),
			TypeRepost, 4,
		},
		{
			regexp.MustCompile("^(https?://)?(m|web)\\.(jellow\\.club|okjike\\.com)/officialMessages?/([0-9a-z]+)(\\?.*)?$"),
			TypeOfficialMessage, 4,
		},
	}
	for _, p := range processors {
		if p.matcher.MatchString(url) {
			return &URL{
				Type: p.postType,
				ID:   p.matcher.FindStringSubmatch(url)[p.idIndex],
			}
		}
	}
	return nil
}

func (u *URL) GenerateMessageUrl() string {
	if u.Type == TypeOfficialMessage {
		return "https://m.okjike.com/officialMessages/" + u.ID
	}
	if u.Type == TypeOriginalPost {
		return "https://m.okjike.com/originalPosts/" + u.ID
	}
	if u.Type == TypeRepost {
		return "https://m.okjike.com/reposts/" + u.ID
	}
	return ""
}
