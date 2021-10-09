package storage

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var longText = strings.Repeat("abc ", 100)

func TestStorage(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tempDir, err := ioutil.TempDir(wd, "TestStorage*")
	t.Logf("temporary directory is %s", tempDir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	st := Storage{
		Path: tempDir,
	}
	err = st.StoreFeedBody([]byte("hello world"))
	assert.NoError(t, err)
	err = st.StoreFeedBody([]byte("foo"))
	assert.NoError(t, err)
	err = st.StoreFeedBody([]byte("hello world"))
	assert.NoError(t, err)
	err = st.StoreFeedBody([]byte(longText))
	assert.NoError(t, err)

	paths, err := st.ListFeedFiles()
	assert.NoError(t, err)

	expected := make(map[string]bool)
	expected["hello world"] = true
	expected["foo"] = true
	expected[longText] = true
	assert.Len(t, paths, len(expected))
	for _, filePath := range paths {
		assert.Len(t, path.Base(filePath), 32+5)
		f, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Error(err)
			continue
		}
		k := string(f)
		if _, ok := expected[k]; !ok {
			t.Errorf("unexpected content: %s", k)
			continue
		}
		delete(expected, k)
	}
	assert.Len(t, expected, 0, "%v", expected)
}
