package service

import (
	"GoNewsScrapper/internals/crawler"
	"GoNewsScrapper/internals/database"
	"GoNewsScrapper/internals/repository"
	"GoNewsScrapper/internals/utils"
	"context"
	"log"
)

type News struct {
	repo *repository.News
}

func NewNews(repo *repository.News) *News {
	return &News{repo: repo}
}

func (n *News) CreateNews(ctx context.Context, news []database.CreateNewsParams) error {
	return n.repo.InsertNewsBulk(ctx, news)
}

func (n *News) GetAllNews(ctx context.Context) ([]database.News, error) {
	return n.repo.GetAllNews(ctx)
}

func (n *News) CrawlNews() error {
	ctx := context.Background()
	news, err := crawler.CrawlProPakistaniPk(ctx, 30)
	if err != nil {
		return err
	}
	news = utils.RemoveNewsDuplicatesBySlug(news)

	err = n.CreateNews(ctx, news)
	if err != nil {
		return err
	}
	log.Println("data insertion complete")

	return nil
}
