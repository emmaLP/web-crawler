package crawler_test

import (
	"io"
	"strings"
	"testing"

	"github.com/emmalp/web-crawler/internal/crawler"
	"github.com/stretchr/testify/assert"
)

func Test_ExtractSameHostLinksFromHTML(t *testing.T) {
	tests := map[string]struct {
		body          io.Reader
		targetHost    string
		expectedLinks []string
	}{
		"empty reader": {
			body:       strings.NewReader(""),
			targetHost: "example.com",
		},
		"invalid html": {
			body:       strings.NewReader("<html>hello/"),
			targetHost: "example.com",
		},
		"a tag without href": {
			body:       strings.NewReader("<html><a/></html>"),
			targetHost: "example.com",
		},
		"single a tag with href": {
			body:          strings.NewReader(`<html><a href="http://example.com"></a></html>`),
			expectedLinks: []string{"http://example.com"},
			targetHost:    "example.com",
		},
		"multiple a tags - only match target host": {
			body:          strings.NewReader(`<html><a href="http://example.com"></a><a href="http://google.com"></a></html>`),
			expectedLinks: []string{"http://example.com"},
			targetHost:    "example.com",
		},
		"self closing href": {
			body:          strings.NewReader(`<html><a href="http://example.com"/></html>`),
			expectedLinks: []string{"http://example.com"},
			targetHost:    "example.com",
		},
		"duplicate link- should only retrun once": {
			body:          strings.NewReader(`<html><a href="http://example.com"/><body><a href="http://example.com"/></body></html>`),
			expectedLinks: []string{"http://example.com"},
			targetHost:    "example.com",
		},
		"nested links": {
			body:          strings.NewReader(`<html><a href="http://example.com"/><body><a href="http://example.com/test"/></body></html>`),
			expectedLinks: []string{"http://example.com", "http://example.com/test"},
			targetHost:    "example.com",
		},
		"no href links": {
			body:       strings.NewReader("<html><body><h1>Hello!</h1></body></html>"),
			targetHost: "example.com",
		},
	}
	for name, cfg := range tests {
		t.Run(name, func(t *testing.T) {
			links := crawler.ExtractSameHostLinksFromHTML(cfg.body, cfg.targetHost)
			assert.ElementsMatch(t, links, cfg.expectedLinks)
		})
	}
}
