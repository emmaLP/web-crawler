package crawler

import (
	"io"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExtractSameHostLinksFromHTML(t *testing.T) {
	testURL, err := url.Parse("http://example.com")
	assert.NoError(t, err)
	tests := map[string]struct {
		body          io.Reader
		targetURL    *url.URL
		expectedLinks []string
	}{
		"empty reader": {
			body:       strings.NewReader(""),
			targetURL: testURL,
		},
		"invalid html": {
			body:       strings.NewReader("<html>hello/"),
			targetURL: testURL,
		},
		"a tag without href": {
			body:       strings.NewReader("<html><a/></html>"),
			targetURL: testURL,
		},
		"single a tag with href": {
			body:          strings.NewReader(`<html><a href="http://example.com"></a></html>`),
			expectedLinks: []string{"http://example.com"},
			targetURL:    testURL,
		},
		"multiple a tags - only match target host": {
			body:          strings.NewReader(`<html><a href="http://example.com"></a><a href="http://google.com"></a></html>`),
			expectedLinks: []string{"http://example.com"},
			targetURL:    testURL,
		},
		"self closing href": {
			body:          strings.NewReader(`<html><a href="http://example.com"/></html>`),
			expectedLinks: []string{"http://example.com"},
			targetURL:    testURL,
		},
		"relative link": {
			body:          strings.NewReader(`<html><a href="/about-us"/></html>`),
			expectedLinks: []string{"http://example.com/about-us"},
			targetURL:    testURL,
		},
		"duplicate link- should only retrun once": {
			body:          strings.NewReader(`<html><a href="http://example.com"/><body><a href="http://example.com"/></body></html>`),
			expectedLinks: []string{"http://example.com"},
			targetURL:    testURL,
		},
		"nested links": {
			body:          strings.NewReader(`<html><a href="http://example.com"/><body><a href="http://example.com/test"/></body></html>`),
			expectedLinks: []string{"http://example.com", "http://example.com/test"},
			targetURL:    testURL,
		},
		"no href links": {
			body:       strings.NewReader("<html><body><h1>Hello!</h1></body></html>"),
			targetURL: testURL,
		},
	}
	for name, cfg := range tests {
		t.Run(name, func(t *testing.T) {
			links := extractSameHostLinksFromHTML(cfg.body, cfg.targetURL)
			assert.ElementsMatch(t, links, cfg.expectedLinks)
		})
	}
}
