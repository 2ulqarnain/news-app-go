package app

import (
	"GoNewsScrapper/internals/database"
	"GoNewsScrapper/internals/repository"
	"GoNewsScrapper/internals/service"
	"context"
	"log"
)

type App struct {
	NewsService *service.News
	Close       func() error
}

func New() (*App, error) {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	queries, err := database.Prepare(context.Background(), db)
	if err != nil {
		log.Fatalln("Error preparing query: ", err)
	}
	err = InitializeDB(queries)
	if err != nil {
		log.Println("Error creating news table ", err)
	}
	repo := repository.NewNews(queries, db)
	svc := service.NewNews(repo)

	return &App{
		NewsService: svc,
		Close: func() error {
			queries.Close()
			return db.Close()
		},
	}, nil
}

func InitializeDB(queries *database.Queries) error {
	return queries.CreateTableNews(context.Background())
}
