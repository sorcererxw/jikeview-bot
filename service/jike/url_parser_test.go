package jike

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUrl(t *testing.T) {
	cases := []struct {
		url    string
		except *Url
	}{
		{
			"https://web.okjike.com/originalPost/5f2f69fac9c6141218b3d73c",
			&Url{OriginalPost, "5f2f69fac9c6141218b3d73c"},
		},
		{
			"https://m.okjike.com/originalPosts/5f2f68947f676b001871a594?s=ewoidSI6ICI1NjQ3MDYyM2U3MjZmNDEyMMjYiCn0=",
			&Url{Type: OriginalPost, ID: "5f2f68947f676b001871a594"},
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.except, ParseUrl(c.url))
	}
}
