package cache

import (
	"bulletin/feed"
	"bulletin/log"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
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

// StoreArticle should not be run in parallel for the same feed, because the writes to the same file might collide and
// garble the content.
func (c *Cache) StoreArticle(article feed.Article) error {
	fname := jsonFileName(article)
	fpath := path.Join(c.cacheDir, fname)
	log.Debugf("cache: store article to %s", fpath)
	fh, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer fh.Close()
	bytes, err := article.Marshall()
	if err != nil {
		return err
	}
	n, err := fh.Write(bytes)
	log.Debugf("cache: wrote %d bytes", n)
	return err
}

func (c *Cache) GetArticles() ([]feed.Article, error) {
	infos, err := ioutil.ReadDir(c.cacheDir)
	if err != nil {
		return nil, err
	}
	var articles []feed.Article
	for _, info := range infos {
		fPath := path.Join(c.cacheDir, info.Name())
		article, err := readArticle(fPath)
		if err != nil {
			log.Debugf("could not read Article from %s: %s", fPath, err)
			continue
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func jsonFileName(art feed.Article) string {
	feedPart := strings.ToLower(escape(art.Feed.Id))
	timePart := art.Published.UTC().Format(time.RFC3339)
	articlePart := escape(art.Id)
	return fmt.Sprintf("%s::%s::%s.json", feedPart, timePart, articlePart)
}

var regexEscape = regexp.MustCompile("[^a-zA-Z0-9]+")

func escape(s string) string {
	return regexEscape.ReplaceAllString(s, "_")
}

func readArticle(path string) (feed.Article, error) {
	fh, err := os.Open(path)
	if err != nil {
		return feed.Article{}, err
	}
	defer fh.Close()
	bytes, err := ioutil.ReadAll(fh)
	if err != nil {
		return feed.Article{}, err
	}
	return feed.UnmarshallArticle(bytes)
}
