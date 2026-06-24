package main

import (
	"GoNewsScrapper/internals/config"
	"GoNewsScrapper/internals/database"
	"GoNewsScrapper/internals/handlers"
	"GoNewsScrapper/internals/repository"
	"GoNewsScrapper/internals/routes"
	"GoNewsScrapper/internals/service"
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	queries, err := database.Prepare(context.Background(), db)
	if err != nil {
		log.Fatalln("Error preparing query: ", err)
	}
	defer queries.Close()
	repo := repository.NewNews(queries, db)
	svc := service.NewNews(repo)
	handler := handlers.NewNews(svc)
	r := routes.NewsRouter(handler)

	addr := fmt.Sprintf(":%s", config.Port)
	fmt.Printf("\nListening on %s...\n", addr)
	http.ListenAndServe(addr, r)
}
