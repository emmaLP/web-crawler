package crawler

import (
	"io"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockHTTPClient mocks the http wrapped client
type mockHTTPClient struct {
	mock.Mock
}

func (m *mockHTTPClient) FetchPage(url string) (io.Reader, error) {
	args := m.Called(url)
	body, ok := args.Get(0).(io.Reader)
	if !ok {
		return nil, args.Error(1)
	}
	return body, args.Error(1)
}

func Test_GenerateSiteMap(t *testing.T) {
	testURL, err := url.Parse("http://example.com")
	require.NoError(t, err)
	tests := map[string]struct {
		expectMocks   func(t *testing.T, httpClient *mockHTTPClient)
		expectedLinks []string
	}{
		"page doesn't contain any links": {
			expectedLinks: []string{"http://example.com"},
			expectMocks: func(t *testing.T, httpClient *mockHTTPClient) {
				t.Helper()
				httpClient.On("FetchPage", "http://example.com").Return(
					io.NopCloser(strings.NewReader("<html>hello/</html>")),
					nil,
				)
			},
		},
		"2 subpages pages with no additional links": {
			expectedLinks: []string{"http://example.com", "http://example.com/about", "http://example.com/home"},
			expectMocks: func(t *testing.T, httpClient *mockHTTPClient) {
				t.Helper()
				httpClient.On("FetchPage", "http://example.com").Return(
					io.NopCloser(strings.NewReader(`<html><a href="http://example.com/about"><a href="http://example.com/home"></html>`)),
					nil,
				)
				httpClient.On("FetchPage", "http://example.com/about").Return(
					io.NopCloser(strings.NewReader(`<html></html>`)),
					nil,
				)
				httpClient.On("FetchPage", "http://example.com/home").Return(
					io.NopCloser(strings.NewReader(`<html></html>`)),
					nil,
				)
			},
		},
		"2 subpages pages with duplicate links": {
			expectedLinks: []string{"http://example.com", "http://example.com/about", "http://example.com/home"},
			expectMocks: func(t *testing.T, httpClient *mockHTTPClient) {
				t.Helper()
				httpClient.On("FetchPage", "http://example.com").Return(
					io.NopCloser(strings.NewReader(`<html><a href="http://example.com/about"><a href="http://example.com/home"></html>`)),
					nil,
				)
				httpClient.On("FetchPage", "http://example.com/about").Return(
					io.NopCloser(strings.NewReader(`<html><a href="http://example.com/home"></html>`)),
					nil,
				)
				httpClient.On("FetchPage", "http://example.com/home").Return(
					io.NopCloser(strings.NewReader(`<html></html>`)),
					nil,
				)
			},
		},
		"2 subpages pages with nested page eache": {
			expectedLinks: []string{"http://example.com",
				"http://example.com/about",
				"http://example.com/home",
				"http://example.com/about/founder",
				"http://example.com/test"},
			expectMocks: func(t *testing.T, httpClient *mockHTTPClient) {
				t.Helper()
				httpClient.On("FetchPage", "http://example.com").Return(
					io.NopCloser(strings.NewReader(`<html><a href="http://example.com/about"><a href="http://example.com/home"></html>`)),
					nil,
				)
				httpClient.On("FetchPage", "http://example.com/about").Return(
					io.NopCloser(strings.NewReader(`<html><a href="http://example.com/about/founder"></html>`)),
					nil,
				)
				httpClient.On("FetchPage", "http://example.com/about/founder").Return(
					io.NopCloser(strings.NewReader(`<html></html>`)),
					nil,
				)
				httpClient.On("FetchPage", "http://example.com/home").Return(
					io.NopCloser(strings.NewReader(`<html><a href="http://example.com/test"><</html>`)),
					nil,
				)
				httpClient.On("FetchPage", "http://example.com/test").Return(
					io.NopCloser(strings.NewReader(`<html></html>`)),
					nil,
				)
			},
		},
	}
	for name, cfg := range tests {
		t.Run(name, func(t *testing.T) {
			httpClient := &mockHTTPClient{}
			if cfg.expectMocks != nil {
				cfg.expectMocks(t, httpClient)
				defer httpClient.AssertExpectations(t)
			}

			client := &Client{
				baseURL:        testURL,
				maxConcurrency: 2,
				httpClient:     httpClient,
			}

			outputChan := make(chan string)
			var outputLinks []string

			go func() {
				for {
					link, more := <-outputChan

					if !more {
						break
					}
					outputLinks = append(outputLinks, link)
				}
			}()

			client.GenerateSiteMap(outputChan)
			close(outputChan)

			assert.ElementsMatch(t, cfg.expectedLinks, outputLinks)
		})
	}
}
