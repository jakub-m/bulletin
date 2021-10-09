package fetcher

import (
	"bulletin/log"
	"sync"
)

const workerCount = 5

type GetAllResult struct {
	Url  string
	Body []byte
	Err  error
}

// GetAll downloads content from multiple URLs.
func GetAll(urls []string) []GetAllResult {
	wg := new(sync.WaitGroup)
	urlChan := make(chan string)
	resultChan := make(chan GetAllResult)
	var results []GetAllResult

	wg.Add(workerCount)
	log.Debugf("use %d workers", workerCount)
	go urlSourcer(urls, urlChan)
	for i := 0; i < workerCount; i++ {
		go worker(wg, urlChan, resultChan)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for r := range resultChan {
		results = append(results, r)
	}
	return results
}

func urlSourcer(urls []string, urlChan chan<- string) {
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)
}

func worker(wg *sync.WaitGroup, urlChan <-chan string, resultChan chan<- GetAllResult) {
	for url := range urlChan {
		log.Infof("Fetch feed from %s", url)
		b, err := Get(url)
		if err != nil {
			log.Infof("Error. Could not fetch %s: %s", url, err)
		}
		log.Debugf("send result to resultChan %dB, %s", len(b), err)
		resultChan <- GetAllResult{
			Url:  url,
			Body: b,
			Err:  err,
		}
	}
	log.Debugf("worker done")
	wg.Done()
}
