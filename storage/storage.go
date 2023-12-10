package storage

import (
	"bulletin/log"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	feedSuffix     = ".feed"
	feedMetaSuffix = ".meta"
)

type Storage struct {
	// Path is the path on the disk where the files will be stored.
	Path string
}

// FeedMeta is ad-hoc solution to address lack of url for Monzo blog. It
// can be simplified by uniting "meta" and "feed body" entities.
type FeedMeta struct {
	Url string `json:"url"`
}

func (st *Storage) StoreFeedBodyMeta(body []byte, url string) error {
	meta := FeedMeta{
		Url: url,
	}
	fname := getFileName(body)
	metaPath := path.Join(st.Path, fname+feedMetaSuffix)
	if err := st.storeFeedMeta(metaPath, meta); err != nil {
		return err
	}
	bodyPath := path.Join(st.Path, fname)
	if err := st.storeFeedBody(bodyPath, body); err != nil {
		return err
	}
	return nil
}

// StoreFeedBody stores body of the feed (raw XML) to a file with name being a hash of its content.
// DEPRECATED
func (st *Storage) StoreFeedBody(body []byte) error {
	fname := getFileName(body)
	bodyPath := path.Join(st.Path, fname)
	return st.storeFeedBody(bodyPath, body)
}

func (st *Storage) storeFeedBody(filePath string, body []byte) error {
	log.Debugf("Write %d B to %s", len(body), filePath)
	return os.WriteFile(filePath, body, 0644)
}

func (st *Storage) storeFeedMeta(filePath string, meta FeedMeta) error {
	j, err := json.MarshalIndent(meta, "", " ")
	if err != nil {
		return err
	}
	log.Debugf("Write %d B to %s", len(filePath), filePath)
	return os.WriteFile(filePath, j, 0644)
}

// ListFeedFiles returns paths to all the stored feed files.
func (st *Storage) ListFeedFiles() ([]string, error) {
	basePath := st.Path
	log.Debugf("List feed files from %s", basePath)
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	feedPaths := []string{}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), feedSuffix) {
			feedPaths = append(feedPaths, path.Join(basePath, e.Name()))
		}
	}
	log.Debugf("Got %d feed paths", len(feedPaths))
	return feedPaths, nil
}

func GetMetaForFeedPath(feedPath string) (FeedMeta, error) {
	metaPath := GetMetaPath(feedPath)
	b, err := ioutil.ReadFile(metaPath)
	if err != nil {
		return FeedMeta{}, err
	}
	meta := &FeedMeta{}
	err = json.Unmarshal(b, meta)
	return *meta, err
}

func GetMetaPath(feedPath string) string {
	return feedPath + feedMetaSuffix
}

func getFileName(body []byte) string {
	return fmt.Sprintf("%x%s", md5.Sum(body), feedSuffix)
}
