package cache

import (
	"encoding/json"
	"feedsummary/feed"
	"feedsummary/log"
	"fmt"
	"os"
	"path"
	"regexp"
	"time"
)

type Cache struct {
	// CacheDir is the directory where the articles are stored.
	cacheDir string
}

func NewCache(cacheDir string) (*Cache, error) {
	c := &Cache{
		cacheDir: cacheDir,
	}
	err := os.MkdirAll(c.cacheDir, 0755)
	return c, err
}

func (c *Cache) StoreArticle(article feed.Article) error {
	fname := jsonFileName(article)
	fpath := path.Join(c.cacheDir, fname)
	log.Debugf("cache: store article to %s", fpath)
	fh, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer fh.Close()
	bytes, err := json.MarshalIndent(article, "", " ")
	if err != nil {
		return err
	}
	n, err := fh.Write(bytes)
	log.Debugf("cache: wrote %d bytes", n)
	return err
}

func jsonFileName(art feed.Article) string {
	return fmt.Sprintf("%s::%s.json", art.Updated.UTC().Format(time.RFC3339), escape(art.Id))
}

var regexEscape = regexp.MustCompile("[^a-zA-Z0-9]+")

func escape(s string) string {
	return regexEscape.ReplaceAllString(s, "_")
}
