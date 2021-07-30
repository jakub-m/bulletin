package compose

import (
	"feedsummary/cache"
	"feedsummary/feed"
	"feedsummary/log"
	"flag"
	"fmt"
	corelog "log"
	"time"
)

const CommandName = "compose"

type Command struct {
	Cache *cache.Cache
}

var referenceTime time.Time

func init() {
	t, err := time.Parse(time.RFC3339, "2000-01-03T00:00:00Z") // monday
	if err != nil {
		corelog.Fatal(err)
	}
	referenceTime = t
}

func (c *Command) Execute(args []string) error {
	opts, err := getOptions(args)
	if err != nil {
		return fmt.Errorf("compose: %s", err)
	}
	intervalStart, err := getNearestInterval(referenceTime, opts.interval, time.Now())
	if err != nil {
		return err
	}
	articles, err := c.Cache.GetArticles()
	if err != nil {
		return err
	}
	var filteredArticles []feed.Article
	for _, a := range articles {
		if a.Updated.After(intervalStart) && !a.Updated.After(intervalStart.Add(opts.interval)) {
			log.Debugf("Accept %s, %s", a.Id, a.Updated)
			filteredArticles = append(filteredArticles, a)
		} else {
			log.Debugf("Drop %s, %s", a.Id, a.Updated)
		}
	}
	formatted, err := feed.FormatHtml(filteredArticles)
	if err != nil {
		return err
	}
	fmt.Println(formatted)
	return nil
}

func getOptions(args []string) (options, error) {
	var options options
	fs := flag.NewFlagSet(CommandName, flag.ContinueOnError)
	fs.DurationVar(&options.interval, "interval", 24*7*time.Hour, "time range of the articles")
	err := fs.Parse(args)
	return options, err
}

type options struct {
	interval time.Duration
}

func getNearestInterval(reference time.Time, interval time.Duration, now time.Time) (time.Time, error) {
	if now.Before(reference) {
		return time.Time{}, fmt.Errorf("compose: cannot find interval %s before %s", now, reference)
	}
	prevPrevInterval, prevInterval := reference, reference
	for t := reference; t.Before(now); t = t.Add(interval) {
		prevPrevInterval, prevInterval = prevInterval, t
	}
	if prevPrevInterval.Add(interval).After(now) {
		return time.Time{}, fmt.Errorf("compose: bug. bad interval end after now: %s", prevPrevInterval.Add(interval))
	}
	return prevPrevInterval, nil
}
