package crawler

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/emmalp/web-crawler/internal/httpwrapper"
)

type httpClient interface {
	FetchPage(url string) (io.Reader, error)
}

// Client contains the necessary infomation to concurrently scrape the given URL
type Client struct {
	baseURL        *url.URL
	maxConcurrency int
	httpClient     httpClient
}

// New instantiates the crawler what will generate asitemap of matching host URLs
// If concurrency is less than 1, the default will be set to 5
func New(url *url.URL, concurrency int) (*Client, error) {
	if url == nil {
		return nil, errors.New("url is required")
	}

	maxConcurrency := 5
	if concurrency > 0 {
		maxConcurrency = concurrency
	}

	return &Client{
		baseURL:        url,
		maxConcurrency: maxConcurrency,
		httpClient:     httpwrapper.New(&http.Client{}),
	}, nil
}

// crawlProcess struct is used internally to this package to keep state per crawl
type crawlProcess struct {
	visited sync.Map
}

// GenerateSiteMap concurrently orchestrates the crawling of the URL and write resulting URL to a channel
// The passed in channel enables the results to be streamed in near real time to the calling code to handle what to do with the information as it receives it.
//
// This func will only look for URLs that match the same domain as the one in the client. This means it will ignore any subdomain.
func (c *Client) GenerateSiteMap(output chan<- string) {
	linksToFetch := make(chan string)
	wg := &sync.WaitGroup{}
	crawler := &crawlProcess{}

	for range c.maxConcurrency {
		wg.Add(1)
		go c.getLinksFromPage(crawler, linksToFetch, output, wg)
	}

	linksToFetch <- c.baseURL.String()

	wg.Wait()
	close(linksToFetch)

}

// getLinksFromPage will continually listen for new links to fetch and determine if the link has been visited or not before outputting it
func (c *Client) getLinksFromPage(crawler *crawlProcess, linksToFetch chan string, output chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	waitCounter := 0
	for {
		select {
		case linkToFetch, ok := <-linksToFetch:
			if !ok {
				return
			}

			// if visited then don't visit again
			if crawler.HaveVisted(linkToFetch) {
				continue
			}

			body, err := c.httpClient.FetchPage(linkToFetch)
			if err != nil {
				log.Printf("unable to fetch %v", linkToFetch)
				continue
			}

			output <- linkToFetch

			links := extractSameHostLinksFromHTML(body, c.baseURL)
			for _, link := range links {
				go func() {
					linksToFetch <- link
				}()
			}
		default:
			// As we read and write from the linksToFetch channel,
			// if we stop receiving messages wait 3 iterations before exiting the func
			if waitCounter > 2 {
				return
			}
			waitCounter++
			time.Sleep(200 * time.Millisecond)

		}
	}

}

// HaveVisted will look to see if the URL has already been visited
// If the record has not been marked as visited then the passed in URL will be stored as visited
func (c *crawlProcess) HaveVisted(url string) bool {
	_, ok := c.visited.Load(url)
	if !ok {
		c.visited.Store(url, 1)
	}

	return ok
}
