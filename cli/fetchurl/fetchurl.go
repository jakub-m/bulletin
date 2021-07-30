package main

import (
	"feedsummary/atom"
	"feedsummary/feed"
	"feedsummary/fetcher"
	"fmt"
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

	a, err := atom.Parse(b)
	if err != nil {
		log.Fatalf("fetchurl: %s", err)
	}
	log.Printf("Got %d articles", len(a.Entries))
	for _, entry := range a.Entries {
		log.Printf("%s\t%s", entry.Uid(), entry.Title)
	}

	html, err := feed.Html(a.GetArticles())
	if err != nil {
		log.Fatalf("fetchurl: %s", err)
	}
	fmt.Println(html)
}
