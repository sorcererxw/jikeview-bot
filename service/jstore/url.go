package jstore

import (
	"fmt"
	"regexp"
)

type (
	Url struct {
		Type string
		ID   string
	}
)

// ParseUrl parses jstore urls
func ParseUrl(url string) *Url {
	re := regexp.MustCompile("^(https?://)?(web|m)\\.jstore\\.site/product/([0-9a-z]+)/?(\\?.*)?$")
	if re.MatchString(url) {
		return &Url{"product", re.FindStringSubmatch(url)[3]}
	}
	return nil
}

func (p Product) generateUrl() string {
	return fmt.Sprintf("https://m.jstore.site/product/%s", p.ID)
}
