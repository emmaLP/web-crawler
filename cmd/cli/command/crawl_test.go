package command

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegration_CrawlRun(t *testing.T) {
	tests := map[string]struct {
		expectedErrMessage string
		url                *url.URL
	}{
		"blank url": {
			url:                nil,
			expectedErrMessage: "url needed to crawl",
		},
		"host is blank - fails to parse url": {
			url: &url.URL{
				Path: `user:abc{DEf1=ghi@example.com:\x7f/hello`,
			},
			expectedErrMessage: "invalid url",
		},
		"host is blank": {
			url: &url.URL{
				Path: "example.com/hello",
			},
		},
		"valid url":{
			url: &url.URL{
				Scheme: "https",
				Host: "example.com",
			},
		},
	}
	for name, cfg := range tests {
		t.Run(name, func(t *testing.T) {

			crawl := &Crawl{cfg.url}

			outputErr := crawl.Run()
			if cfg.expectedErrMessage != "" {
				assert.ErrorContains(t, outputErr, cfg.expectedErrMessage)
			} else {
				assert.NoError(t, outputErr)
			}
		})
	}
}
