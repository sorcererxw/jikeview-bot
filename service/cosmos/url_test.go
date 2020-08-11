package cosmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUrl(t *testing.T) {
	cases := []struct {
		Url    string
		Except interface{}
	}{
		{
			"https://www.xiaoyuzhoufm.com/episode/5f1aab919777af96d?s=eyJ1IjoiNiYzUxYzE5MGY2IiwiZCI6M30%3D",
			&Url{ID: "5f1aab919777af96d"},
		},
		{
			"https://www.xiaoyuzhoufm.com/episode/5f1aab919777af96d/comment",
			nil,
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.Except == nil, ParseUrl(c.Url) == nil)
		if c.Except != nil {
			assert.Equal(t, c.Except, ParseUrl(c.Url))
		}
	}

}
