package crawler

import (
	"io"
	"net/url"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

// ExtractSameHostLinksFromHTML will take the html body within the reader and find any a tag that have an href attribute
// This function will all ensure that the links returned from the page only match the specified host
// Any links from external sources or subdomain on the request host will be ignored.
func extractSameHostLinksFromHTML(body io.Reader, targetURL *url.URL) []string {
	var links []string

	htmTokenizer := html.NewTokenizer(body)
	for {
		currentPosition := htmTokenizer.Next()
		//Check if reached the end of the input or an error has occurred
		if currentPosition == html.ErrorToken {
			log.Debug().Msgf("end of input")
			// Sorting the URLs and then compacting them to remove any duplicate links
			slices.Sort(links)
			return slices.Compact(links)
		}

		token := htmTokenizer.Token()
		// Find any start tag or a self closing tag and only if if the token is `a`
		if (currentPosition == html.StartTagToken || currentPosition == html.SelfClosingTagToken) && token.Data == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					link := parseReturnedLink(targetURL, attr.Val)
					if isSameHostName(link, targetURL.Host) {
						links = append(links, link)
						break
					}
				}
			}
		}
	}

}

func parseReturnedLink(base *url.URL, link string) string {
	log.Print("Link:", link)
	newLink := link
	if strings.HasPrefix(link, "/") {
		log.Print("Has prefix")
		newLink = base.JoinPath(link).String()
		log.Print("New link:", newLink)
	}

	return strings.TrimSuffix(newLink, "/")
}

func isSameHostName(foundURL, targetHost string) bool {
	found, err := url.Parse(foundURL)
	if err != nil {
		return false
	}

	return found.Host == targetHost
}
