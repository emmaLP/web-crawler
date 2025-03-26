package crawler_test

import (
	"io"
	"strings"
	"testing"

	"github.com/emmalp/web-crawler/internal/crawler"
	"github.com/stretchr/testify/assert"
)

func Test_ExtractLinksFromHTML(t *testing.T){
	tests := map[string]struct{
		body io.Reader
		expectedLinks []string
	}{
		"empty reader":{
			body:  strings.NewReader(""),
		},
		"invalid html":{
			body:  strings.NewReader("<html>hello/"),
		},
		"a tag without href":{
			body:  strings.NewReader("<html><a/></html>"),
		},
		"single a tag with href":{
			body:  strings.NewReader(`<html><a href="http://example.com"></a></html>`),
			expectedLinks: []string{"http://example.com"},
		},
		"multiple a tags":{
			body:  strings.NewReader(`<html><a href="http://example.com"></a><a href="http://google.com"></a></html>`),
			expectedLinks: []string{"http://google.com","http://example.com"},
		},
		"self closing href":{
			body:  strings.NewReader(`<html><a href="http://example.com"/></html>`),
			expectedLinks: []string{"http://example.com"},
		},
		"duplicate link- should only retrun once":{
			body:  strings.NewReader(`<html><a href="http://example.com"/><body><a href="http://example.com"/></body></html>`),
			expectedLinks: []string{"http://example.com"},
		},
		"nested links":{
			body:  strings.NewReader(`<html><a href="http://example.com"/><body><a href="http://google.com"/></body></html>`),
			expectedLinks: []string{"http://example.com", "http://google.com"},
		},
		"no href links":{
			body:  strings.NewReader("<html><body><h1>Hello!</h1></body></html>"),
		},
	}
	for name, cfg := range tests {
		t.Run(name, func(t *testing.T) {
			links := crawler.ExtractLinksFromHTML(cfg.body)
			assert.ElementsMatch(t, links, cfg.expectedLinks)
		})
	}
}
