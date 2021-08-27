package storage

import (
	"bulletin/log"
	"crypto/md5"
	"fmt"
	"os"
	"path"
	"strings"
)

const feedSuffix = ".feed"

type Storage struct {
	// Path is the path on the disk where the files will be stored.
	Path string
}

// StoreFeedBody stores body of the feed (raw XML) to a file with name being a hash of its content.
func (st *Storage) StoreFeedBody(body []byte) error {
	fname := getFileName(body)
	p := path.Join(st.Path, fname)
	log.Debugf("Write %d B to %s", len(body), p)
	return os.WriteFile(p, body, 0644)
}

// ListFiles returns paths to all the stored feed files.
func (st *Storage) ListFiles() ([]string, error) {
	basePath := st.Path
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
	return feedPaths, nil
}

func getFileName(body []byte) string {
	return fmt.Sprintf("%x%s", md5.Sum(body), feedSuffix)
}
