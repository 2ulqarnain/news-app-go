package main

import (
	newsApp "GoNewsScrapper/internals/app"
	"log"
)

func main() {
	app, err := newsApp.New()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	if err := app.NewsService.CrawlNews(); err != nil {
		log.Fatal(err)
	}
}
