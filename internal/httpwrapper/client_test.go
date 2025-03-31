package httpwrapper

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockHTTPClient struct {
	mock.Mock
}

func (m *mockHTTPClient) Get(url string) (resp *http.Response, err error) {
	args := m.Called(url)
	var output *http.Response
	if out, ok := args.Get(0).(*http.Response); ok {
		output = out
	}

	return output, args.Error(1)
}

func Test_FetchPage(t *testing.T) {
	tests := map[string]struct {
		url                string
		expectedMocks      func(t *testing.T, httpClient *mockHTTPClient, url string)
		expectedBody       string
		expectedErrMessage string
	}{
		"successfully fetch page": {
			url: "http://hello.com/why",
			expectedMocks: func(t *testing.T, httpClient *mockHTTPClient, url string) {
				t.Helper()
				httpClient.On("Get", url).Return(&http.Response{
					Body:       io.NopCloser(bytes.NewBufferString(`<html></html>`)),
					StatusCode: http.StatusOK,
				}, nil)
			},
			expectedBody: "<html></html>",
		},
		"error fetching page": {
			url: "http://error.com/why",
			expectedMocks: func(t *testing.T, httpClient *mockHTTPClient, url string) {
				t.Helper()
				httpClient.On("Get", url).Return( nil, assert.AnError)
			},
			expectedErrMessage: "unable to perform get request to url (http://error.com/why)",
		},
	}
	for name, cfg := range tests {
		t.Run(name, func(t *testing.T) {
			mockHTTP := &mockHTTPClient{}
			if cfg.expectedMocks != nil {
				cfg.expectedMocks(t, mockHTTP, cfg.url)
				defer mockHTTP.AssertExpectations(t)
			}

			client := New(mockHTTP)
			body, err := client.FetchPage(cfg.url)
			if cfg.expectedErrMessage != "" {
				assert.ErrorContains(t, err, cfg.expectedErrMessage)
			} else {
				assert.NoError(t, err)

				strBody, err := io.ReadAll(body)
				require.NoError(t, err)
				assert.Equal(t, cfg.expectedBody, string(strBody))
			}
		})
	}
}
