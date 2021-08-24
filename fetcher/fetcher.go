package fetcher

import (
	"bulletin/log"
	"fmt"
	"io"
	"net/http"
)

// Get fetches a raw content of the url.
func Get(url string) ([]byte, error) {
	visited := make(map[string]bool)
	return getRec(url, visited)
}

func getRec(url string, visited map[string]bool) ([]byte, error) {
	log.Debugf("fetcher: get url: %s", url)
	if _, alreadyVisited := visited[url]; alreadyVisited {
		return nil, fetcherError(fmt.Errorf("url already visited, possible infinite redirection: %s", url))
	}
	visited[url] = true

	resp, err := http.Get(url)
	if err != nil {
		return nil, fetcherError(err)
	}
	defer resp.Body.Close()
	if url, err := resp.Location(); err == nil {
		log.Debugf("fetcher: redirection to %s", url)
		return getRec(url.String(), visited)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fetcherError(nil)
	}
	log.Debugf("fetcher: read %dkB", len(body)/1000)
	return body, nil
}

func fetcherError(e error) error {
	return fmt.Errorf("fetcher: %s", e)
}
