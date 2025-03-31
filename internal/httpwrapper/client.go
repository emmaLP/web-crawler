package httpwrapper

import (
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type httpClient interface {
	Get(url string) (resp *http.Response, err error)
}

// Client struct holds a single instance of the httpwrapper client so that we are not instantiating an http client per API call
type Client struct {
	httpClient httpClient
}

// New instantiates http client struct which will either use the passed in httpClient or create a new one
func New(httpClient httpClient) *Client {
	hClient := httpClient
	if httpClient == nil {
		hClient = &http.Client{}
	}
	return &Client{
		httpClient: hClient,
	}
}

// FetchPage fetches the page based on the URL and returns the body in reader format or an error if something goes wrong
// url is a function variable as the url will likely change per request within the same instantiation
func (c *Client) FetchPage(url string) (io.Reader, error) {
	result, err := c.httpClient.Get(url)
	if err != nil{
		return nil, fmt.Errorf("unable to perform get request to url (%v): %w", url, err)
	}
	log.Debug().Msgf("Result status: %d", result.StatusCode)
	return result.Body, nil
}
