package command

import (
	"bufio"
	"bulletin/fetcher"
	"bulletin/log"
	"bulletin/storage"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

const FetchCommandName = "fetch"

var defaultConfigPath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defaultConfigPath = path.Join(home, ".bulletin", "feeds.conf")
}

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
	urls, err := readFeedUrls(opts.feedsConfig)
	if err != nil {
		return err
	}
	for _, u := range urls {
		log.Debugf("url: %s", u)
	}

	for _, r := range fetcher.GetAll(urls) {
		if r.Err == nil {
			err = c.Storage.StoreFeedBody(r.Body)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type fetchOptions struct {
	feedsConfig string
}

func getFetchOptions(args []string) (fetchOptions, error) {
	var options fetchOptions
	fs := flag.NewFlagSet(FetchCommandName, flag.ContinueOnError)
	os.Getwd()
	fs.StringVar(&options.feedsConfig, "feeds", defaultConfigPath, "path to feed configuration. Use `-` for standard input")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage for fetch: pass URLs as positional arguments.\n")
		fs.PrintDefaults()
	}
	err := fs.Parse(args)
	return options, err
}

func readFeedUrls(path string) ([]string, error) {
	var urls []string
	var fh *os.File
	if path == "-" {
		fh = os.Stdin
	} else {
		var err error
		fh, err = os.Open(path)
		if err != nil {
			return nil, err
		}
		defer fh.Close()
	}
	s := bufio.NewScanner(fh)
	for s.Scan() {
		line := s.Text()
		if i := strings.Index(line, "#"); i >= 0 {
			line = line[:i]
		}
		line = strings.Trim(line, " \n\r")
		if line != "" {
			urls = append(urls, line)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}
