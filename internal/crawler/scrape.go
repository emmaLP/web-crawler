package crawler

import (
	"errors"
	"net/http"
	"net/url"
)

// Clientstructu containing the necessary infomation to concurrently scrape the given URL
type Client struct {
	url            *url.URL
	maxConcurrency int
	httpClient     *http.Client
}

// New instantiates the crawler what will generate asitemap of matching host URLs
// If concurrency is less than 1, the default will be set to 5
func New(url *url.URL, concurrency int) (*Client, error) {
	if url == nil {
		return nil, errors.New("url is required")
	}

	maxConcurrency := 5
	if concurrency > 0 {
		maxConcurrency = 5
	}

	return &Client{
		url:            url,
		maxConcurrency: maxConcurrency,
		httpClient:     &http.Client{},
	}, nil
}

 // crawlProcess struct is used internally to this package to keep state per crawl
type crawlProcess struct {
	visited map[string]int
	
}

// GenerateSiteMap concurrently orchestrates the crawling of the URL and write resulting URL to a channel
// The passed in channel enables the results to be streamed in near real time to the calling code to handle what to do with the information as it receives it.
//
// This func will only look for URLs that match the same domain as the one in the client. This means it will ignore any subdomain.
func (c *Client) GenerateSiteMap(output chan<- string) error {

	return nil
}
