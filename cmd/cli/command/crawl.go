package command

// Crawl struct contains the CLI args needed to start the crawl command
type Crawl struct {
	Domain string `short:"d" help:"Domain to web crawl. You don't need to include the scheme (https/http)"`
	Scheme string `short:"s" optional:"" default:"https" help:"Overrides the scheme for the URL used. Defaults to https"`,

}

func (c *Crawl) Run() error {

	return nil
}