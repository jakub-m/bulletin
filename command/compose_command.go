package command

import (
	"bulletin/feed"
	"bulletin/feedparser"
	"bulletin/log"
	"bulletin/storage"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	corelog "log"
	"os"
	"path"
	"sort"
	"time"
)

const (
	ComposeCommandName = "compose"
	durationDay        = time.Hour * 24
	filenameTimeLayout = "2006-01-02"
	referenceMonday    = "2000-01-03T00:00:00Z"
)

type ComposeCommand struct {
	Storage *storage.Storage
}

var referenceTime time.Time

func init() {
	t, err := time.Parse(time.RFC3339, referenceMonday)
	if err != nil {
		corelog.Fatal(err)
	}
	referenceTime = t.Add(4 * durationDay) // friday
}

func (c *ComposeCommand) Execute(args []string) error {
	now := time.Now()
	opts, err := getComposeOptions(args)
	if err != nil {
		return fmt.Errorf("compose: %s", err)
	}
	interval := time.Duration(opts.intervalDays) * 24 * time.Hour
	intervalStart := getNearestInterval(referenceTime, interval, now)
	intervalEnd := intervalStart.Add(interval)
	feeds := c.getFeeds()
	feeds = filterArticlesInFeeds(feeds, intervalStart, intervalEnd)
	sortFeeds(feeds)

	var pageTemplate *string
	if opts.templatePath != "" {
		log.Infof("Use page template: %s", opts.templatePath)
		f, err := ioutil.ReadFile(opts.templatePath)
		if err != nil {
			return err
		}
		t := string(f)
		pageTemplate = &t
	}

	formatted, err := feed.FormatFeedsAsHtml(opts.intervalDays, intervalEnd, pageTemplate, feeds)
	if err != nil {
		return err
	}
	w, actualPath, err := newOutput(opts.output, intervalEnd)
	if err != nil {
		return err
	}
	defer w.Close()
	log.Infof("output: %s", actualPath)
	fmt.Fprintln(w, formatted)
	return nil
}

func (c *ComposeCommand) getFeeds() []feed.Feed {
	feedPaths, err := c.Storage.ListFeedFiles()
	if err != nil {
		log.Infof("Failed to list files: %s", err)
	}
	feeds := []feed.Feed{}
	for _, feedPath := range feedPaths {
		log.Debugf("Parse %s", feedPath)
		body, err := ioutil.ReadFile(feedPath)
		if err != nil {
			log.Infof("Failed to open %s: %v", feedPath, err)
			continue
		}
		meta, err := storage.GetMetaForFeedPath(feedPath)
		if err != nil {
			log.Infof("Failed to load meta for %s: %v", feedPath, err)
			continue
		}
		f, err := feedparser.GetFeed(body, meta.Url)
		feed.FixRelativeUrls(&f)
		if err != nil {
			log.Infof("Failed to parse %s: %s", feedPath, err)
			continue
		}
		feeds = append(feeds, f)
	}
	return feeds
}

func filterArticlesInFeeds(feeds []feed.Feed, intervalStart, intervalEnd time.Time) []feed.Feed {
	var newFeeds []feed.Feed
	for _, f := range feeds {
		var filteredArticles []feed.Article
		for _, a := range f.Articles {
			if a.Published.After(intervalStart) && !a.Published.After(intervalEnd) {
				log.Debugf("Accept %s, %s", a.Id, a.Published)
				filteredArticles = append(filteredArticles, a)
			} else {
				log.Debugf("Drop %s, %s", a.Id, a.Published)
			}
		}
		newFeed := f
		newFeed.Articles = filteredArticles
		if len(newFeed.Articles) > 0 {
			newFeeds = append(newFeeds, newFeed)
		}
	}
	return newFeeds
}

func sortFeeds(feeds []feed.Feed) {
	sort.Slice(feeds, func(i, j int) bool {
		f := feeds[i]
		g := feeds[j]
		if len(f.Articles) != len(g.Articles) {
			return len(f.Articles) < len(g.Articles)
		}
		// shuffles the feeds deterministically
		return feedSortHash(f) < feedSortHash(f)
	})
}

func feedSortHash(f feed.Feed) uint32 {
	hash := fnv.New32()
	fmt.Fprintf(hash, f.Id)
	for _, art := range f.Articles {
		fmt.Fprint(hash, art.Id)
	}
	return hash.Sum32()
}

func newOutput(outPath string, intervalEnd time.Time) (io.WriteCloser, string, error) {
	if outPath == "-" {
		return &nopCloser{os.Stdout}, "stdout", nil
	}
	if fileInfo, err := os.Stat(outPath); err == nil && fileInfo.IsDir() {
		fname := fmt.Sprintf("bulletin-%s.md", intervalEnd.Format(filenameTimeLayout))
		outPath = path.Join(outPath, fname)
	}
	w, err := os.Create(outPath)
	return w, outPath, err
}

func getComposeOptions(args []string) (composeOptions, error) {
	var options composeOptions
	fs := flag.NewFlagSet(ComposeCommandName, flag.ContinueOnError)
	fs.IntVar(&options.intervalDays, "days", 7, "time range of the articles in DAYS")
	fs.StringVar(&options.templatePath, "template", "", "template to render the bulletin")
	fs.StringVar(&options.output, "output", ".", "output. can be directory, concrete file name or `-` for stdout.")
	err := fs.Parse(args)
	return options, err
}

type composeOptions struct {
	intervalDays int
	templatePath string
	output       string
}

func getNearestInterval(reference time.Time, interval time.Duration, now time.Time) time.Time {
	n := now.Sub(reference) / interval
	d := (n - 1) * interval
	return reference.Add(d)
}

type nopCloser struct {
	f *os.File
}

func (c *nopCloser) Close() error {
	return nil
}

func (c *nopCloser) Write(p []byte) (n int, err error) {
	return c.f.Write(p)
}
