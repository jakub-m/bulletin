package main

import (
	"feedsummary/command"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	err := mainerr()
	if err != nil {
		fmt.Printf("ERROR. %s\n", err)
		os.Exit(1)
	}
}

func mainerr() error {
	commands := make(map[string]command.Command)
	commands["fetch"] = &command.FetchCommand{}

	flag.Usage = func() {
		fmt.Printf("Available commands: %s\n", strings.Join(keys(commands), ", "))
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return fmt.Errorf("missing command")
	}
	commandString := flag.Arg(0)
	cmd, ok := commands[commandString]
	if !ok {
		return fmt.Errorf("unknown command: %s", commandString)
	}
	err := cmd.Execute(flag.Args()[1:])
	if err != nil {
		return err
	}
	return nil
}

func keys(m map[string]command.Command) []string {
	var stringKeys []string
	for k := range m {
		stringKeys = append(stringKeys, k)
	}
	return stringKeys
}
