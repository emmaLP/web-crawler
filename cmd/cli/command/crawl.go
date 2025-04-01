package command

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/emmalp/web-crawler/internal/crawler"
)

// Crawl struct contains the CLI args needed to start the crawl command
type Crawl struct {
	URL *url.URL `arg:"" help:"URL of the website to crawl. This can be just the domain or include the scheme. If no scheme is provided, it will default to https. For example: https://monzo.com or monzon.com"`
}

// Run executes the command with the necessary fields.
// This Run func signature is defined by the Kong package to support binidng more than one command in a single CLI
//
// See [here](https://github.com/alecthomas/kong?tab=readme-ov-file#attach-a-run-error-method-to-each-command) for more info
func (c *Crawl) Run() error {
	if c.URL == nil {
		return errors.New("url needed to crawl")
	}

	if c.URL.Scheme == "" {
		c.URL.Scheme = "https"
	}
	
	if c.URL.Host == "" {
		//if you start the CLI with just a domain such as monzo.com, the kong cli parses that as a url path not the host
		// we correct this here by re-running url.Parse now that we have a scheme set
		url, err := url.Parse(c.URL.String())
		if err != nil {
			return fmt.Errorf("invalid url: %w", err)
		}
		c.URL = url
	}
	
	// Blank any path as we only want the domain
	// This ensures no duplicates especially if there is trailing slash
	c.URL.Path = ""

	links := make(chan string)
	crawlerClient, err := crawler.New(c.URL, 5)
	if err != nil {
		return fmt.Errorf("unable to instantiate crawler: %w", err)
	}

	// Log out the links as they are added to the channel
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		counter := 0
		for {
			link, more := <-links

			if !more {
				break
			}
			counter++
			log.Print(link)

		}
		log.Printf("processed links: %d", counter)
	}()

	crawlerClient.GenerateSiteMap(links)

	close(links)
	wg.Wait()

	return nil
}
