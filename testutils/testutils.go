package testutils

import (
	"bulletin/feed"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ParseTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("parseTime: %s", err)
	}
	return parsed
}

func ParseFromFile(t *testing.T, parser feed.FeedParser, filePath, url string) feed.Feed {
	file, err := os.Open(filePath)
	assert.NoError(t, err)
	defer file.Close()
	body, err := io.ReadAll(file)
	assert.NoError(t, err)
	fe, err := parser.ParseFeed(body, url)
	assert.NoError(t, err)
	return fe
}
