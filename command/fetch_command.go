package command

import (
	"bulletin/fetcher"
	"bulletin/log"
	"bulletin/storage"
	"flag"
	"fmt"
	"os"
)

const FetchCommandName = "fetch"

// FetchCommand fetches feed from a single source provided directly in the command line.
type FetchCommand struct {
	Storage *storage.Storage
}

func (c *FetchCommand) Execute(args []string) error {
	opts, err := getFetchOptions(args)
	if err != nil {
		return err
	}
	log.Debugf("options: %+v", opts)
	for _, url := range opts.urls {
		log.Infof("Fetch feed from %s", url)
		feedBody, err := fetcher.Get(url)
		if err != nil {
			log.Infof("Error. Could not fetch %s: %s", url, err)
			continue
		}
		err = c.Storage.StoreFeedBody(feedBody)
		if err != nil {
			return err
		}
	}
	return nil
}

func getFetchOptions(args []string) (fetchOptions, error) {
	var options fetchOptions
	fs := flag.NewFlagSet(FetchCommandName, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage for fetch: pass URLs as positional arguments.\n")
		fs.PrintDefaults()
	}
	err := fs.Parse(args)
	if fs.NArg() == 0 {
		//lint:ignore ST1005 the error is printed with usage and would look weird.
		return options, fmt.Errorf("Missing URLs as positional arguments.")
	}
	options.urls = fs.Args()
	return options, err
}

type fetchOptions struct {
	urls []string
}
