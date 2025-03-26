package command

import (
	"log"
	"net/url"
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
	if c.URL.Scheme == "" {
		c.URL.Scheme= "https"
	}
	log.Print("URL:", c.URL.String())

	//TODO Crawl here
	return nil
}