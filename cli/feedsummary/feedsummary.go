package main

import (
	"feedsummary/cache"
	"feedsummary/command"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

const cacheBaseName = "feedsummary_cache"

func main() {
	err := mainErr()
	if err != nil {
		fmt.Printf("ERROR. %s\n", err)
		os.Exit(1)
	}
}

func mainErr() error {
	availableCommands := []string{command.FetchCommandName, command.ComposeCommandName}

	flag.Usage = func() {
		fmt.Printf("Available commands: %s\n", strings.Join(availableCommands, ", "))
		flag.PrintDefaults()
	}
	var opts options
	defaultCacheDir := path.Join(os.TempDir(), cacheBaseName)
	flag.StringVar(&opts.CacheDir, "cache", defaultCacheDir, "cache directory")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return fmt.Errorf("missing command")
	}

	cacheInstance, err := cache.NewCache(opts.CacheDir)
	if err != nil {
		return err
	}
	commands := make(map[string]command.Command)
	commands[command.FetchCommandName] = &command.FetchCommand{
		Cache: cacheInstance,
	}
	commands[command.ComposeCommandName] = &command.ComposeCommand{
		Cache: cacheInstance,
	}

	commandString := flag.Arg(0)
	cmd, ok := commands[commandString]
	if !ok {
		return fmt.Errorf("unknown command: %s", commandString)
	}
	return cmd.Execute(flag.Args()[1:])
}

type options struct {
	CacheDir string
}
