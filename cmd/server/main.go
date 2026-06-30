package main

import (
	newsApp "GoNewsScrapper/internals/app"
	"GoNewsScrapper/internals/config"
	"GoNewsScrapper/internals/handlers"
	"GoNewsScrapper/internals/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	app, err := newsApp.New()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()
	handler := handlers.NewNews(app.NewsService)
	r := routes.NewsRouter(handler)

	addr := fmt.Sprintf(":%s", config.Port)
	fmt.Printf("\nListening on %s...\n", addr)
	http.ListenAndServe(addr, r)
}
