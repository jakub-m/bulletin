package main

import (
	"bulletin/command"
	"bulletin/log"
	"bulletin/storage"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

const cacheBaseName = "bulletin_cache"

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
	flag.StringVar(&opts.cacheDir, "cache", defaultCacheDir, "cache directory")
	flag.BoolVar(&opts.verbose, "verbose", false, "verbose log")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return fmt.Errorf("missing command")
	}
	log.SetVerbose(opts.verbose)

	storageInstance := &storage.Storage{
		Path: opts.cacheDir,
	}
	commands := make(map[string]command.Command)
	commands[command.FetchCommandName] = &command.FetchCommand{
		Storage: storageInstance,
	}
	commands[command.ComposeCommandName] = &command.ComposeCommand{
		Storage: storageInstance,
	}

	commandString := flag.Arg(0)
	cmd, ok := commands[commandString]
	if !ok {
		return fmt.Errorf("unknown command: %s", commandString)
	}
	return cmd.Execute(flag.Args()[1:])
}

type options struct {
	cacheDir string
	verbose  bool
}
