package command

import "net/url"

// Crawl struct contains the CLI args needed to start the crawl command
type Crawl struct {
	URL *url.URL `short:"u" help:"Scheme and Domain to web crawl."`

}

// Run executes the command with the necessary fields. 
// This Run func signature is defined by the Kong package to support binidng more than one command in a single CLI
// 
// See [here](https://github.com/alecthomas/kong?tab=readme-ov-file#attach-a-run-error-method-to-each-command) for more info
func (c *Crawl) Run() error {

	return nil
}