package main

import (
	"feedsummary/fetcher"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("expected url as an argumetn")
	}
	b, err := fetcher.Get(args[1])
	if err != nil {
		log.Fatalf("fetchurl: %s", err)
	}
	log.Printf("fetchurl: Got %d bytes", len(b))
}
