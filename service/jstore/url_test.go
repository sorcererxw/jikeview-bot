package jstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sorcererxw/jikeview-bot/service/jstore"
)

func TestParseUrl(t *testing.T) {
	cases := [][]string{
		{"https://m.jstore.site/product/123", "123"},
		{"https://web.jstore.site/product/123", "123"},
		{"https://m.jstore.site/product/123?source=tg", "123"},
		{"https://m.jstore.site/product/123/?source=tg", "123"},
	}

	for _, c := range cases {
		assert.Equal(t, jstore.ParseUrl(c[0]).ID, c[1])
	}
}
