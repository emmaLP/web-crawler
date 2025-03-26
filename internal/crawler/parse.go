package crawler

import (
	"io"
	"slices"

	"golang.org/x/net/html"
)

func ExtractLinksFromHTML(body io.Reader)[]string {
	var links []string

	htmTokenizer := html.NewTokenizer(body)
	for {
		currentPosition := htmTokenizer.Next() 
		//Check if reached the end of the input or an error has occurred
		if currentPosition == html.ErrorToken {
			// Sorting the URLs and then compacting them to remove any duplicate links
			slices.Sort(links)
			return slices.Compact(links)
		}


		token := htmTokenizer.Token()
		// Find any start tag or a self closing tag and only if if the token is `a`
		if (currentPosition == html.StartTagToken || currentPosition == html.SelfClosingTagToken )&& token.Data == "a"{
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
						break
					}
			}
		}

	}

}