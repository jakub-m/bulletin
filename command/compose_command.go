package command

import (
	"bulletin/feed"
	"bulletin/feedparser"
	"bulletin/log"
	"bulletin/storage"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	corelog "log"
	"os"
	"path"
	"time"
)

const ComposeCommandName = "compose"
const filenameTimeLayout = "2006-01-02"

type ComposeCommand struct {
	Storage *storage.Storage
}

var referenceTime time.Time

func init() {
	t, err := time.Parse(time.RFC3339, "2000-01-03T00:00:00Z") // monday
	if err != nil {
		corelog.Fatal(err)
	}
	referenceTime = t
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
	articles := c.getArticles()

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

	var filteredArticles []feed.Article
	for _, a := range articles {
		if a.Published.After(intervalStart) && !a.Published.After(intervalEnd) {
			log.Debugf("Accept %s, %s", a.Id, a.Published)
			filteredArticles = append(filteredArticles, a)
		} else {
			log.Debugf("Drop %s, %s", a.Id, a.Published)
		}
	}

	formatted, err := feed.FormatHtml(opts.intervalDays, intervalEnd, pageTemplate, filteredArticles)
	if err != nil {
		return err
	}
	w, err := newOutput(opts.output, intervalEnd)
	if err != nil {
		return err
	}
	defer w.Close()
	fmt.Fprintln(w, formatted)
	return nil
}

func (c *ComposeCommand) getArticles() []feed.Article {
	paths, err := c.Storage.ListFiles()
	articles := []feed.Article{}
	if err != nil {
		log.Infof("Failed to list files: %s", err)
	}
	for _, path := range paths {
		log.Debugf("Parse %s", path)
		b, err := ioutil.ReadFile(path)
		if err != nil {
			log.Infof("Failed to open %s: %s", path, err)
			continue
		}
		a, err := feedparser.GetArticles(b)
		if err != nil {
			log.Infof("Failed to parse %s: %s", path, err)
			continue
		}
		articles = append(articles, a...)
	}
	return articles
}

func newOutput(outPath string, intervalEnd time.Time) (io.WriteCloser, error) {
	if outPath == "-" {
		return &nopCloser{os.Stdout}, nil
	}
	if fileInfo, err := os.Stat(outPath); err == nil && fileInfo.IsDir() {
		fname := fmt.Sprintf("bulletin-%s.html", intervalEnd.Format(filenameTimeLayout))
		outPath = path.Join(outPath, fname)
	}
	return os.Create(outPath)
}

func getComposeOptions(args []string) (composeOptions, error) {
	var options composeOptions
	fs := flag.NewFlagSet(ComposeCommandName, flag.ContinueOnError)
	fs.IntVar(&options.intervalDays, "days", 7, "time range of the articles in DAYS")
	fs.StringVar(&options.templatePath, "template", "", "template to render the bulletin")
	fs.StringVar(&options.output, "output", "-", "output. can be directory, concrete file name or `-` for stdout.")
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
