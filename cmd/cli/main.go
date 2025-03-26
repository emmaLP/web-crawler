package main

import (
	"github.com/alecthomas/kong"
	"github.com/emmalp/web-crawler/cmd/cli/command"
)

var cli struct {
  Crawler command.Crawl `cmd:"" help:""`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("crawler"),
		kong.UsageOnError())

  err := ctx.Run()
  ctx.FatalIfErrorf(err)
}
