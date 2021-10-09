package command

import (
	"bulletin/log"
	"bulletin/storage"
	"fmt"
	"os"
)

const CleanCommandName = "clean"

type CleanCommand struct {
	Storage *storage.Storage
}

func (c *CleanCommand) Execute(args []string) error {
	if len(args) > 0 {
		fmt.Println("clean command cleans the file cache.")
		return nil
	}
	succeed := 0
	feedPaths, err := c.Storage.ListFeedFiles()
	if err != nil {
		return err
	}
	remove := func(path string) {
		if err := os.Remove(path); err == nil {
			succeed++
		} else {
			log.Infof("failed to remove %s: %v", path, err)
		}
	}
	for _, feedPath := range feedPaths {
		remove(feedPath)
		remove(storage.GetMetaPath(feedPath))
	}
	log.Infof("removed %d files from %s", succeed, c.Storage.Path)
	return nil
}
