package jstore

import (
	"fmt"
	"regexp"
)

type (
	URL struct {
		Type string
		ID   string
	}
)

// ParseURL parses jstore urls
func ParseURL(url string) *URL {
	re := regexp.MustCompile(`^(https?://)?(web|m)\.jstore\.site/product/([0-9a-z]+)/?(\?.*)?$`)
	if re.MatchString(url) {
		return &URL{"product", re.FindStringSubmatch(url)[3]}
	}
	return nil
}

func (p Product) generateURL() string {
	return fmt.Sprintf("https://m.jstore.site/product/%s", p.ID)
}
