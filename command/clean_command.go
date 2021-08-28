package command

import (
	"bulletin/log"
	"bulletin/storage"
	"os"
)

const CleanCommandName = "clean"

type CleanCommand struct {
	Storage *storage.Storage
}

func (c *CleanCommand) Execute(args []string) error {
	succeed := 0
	files, err := c.Storage.ListFiles()
	if err != nil {
		return err
	}
	for _, path := range files {
		if err := os.Remove(path); err == nil {
			succeed++
		} else {
			log.Infof("failed to remove %s: %s", path, err)
		}
	}
	log.Infof("removed %d files from %s", succeed, c.Storage.Path)
	return nil
}
