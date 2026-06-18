package service

import (
	"GoNewsScrapper/internals/database"
	"GoNewsScrapper/internals/repository"
	"context"
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
