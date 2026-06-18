package repository

import (
	"GoNewsScrapper/internals/database"
	"context"
	"database/sql"
)

type NewsInterface interface {
	GetAllNews(context.Context) ([]database.News, error)
	InsertNewsBulk(context.Context, []database.CreateNewsParams) error
}

type News struct {
	q  *database.Queries
	db *sql.DB
}

func NewNews(q *database.Queries, db *sql.DB) *News {
	return &News{q: q, db: db}
}

func (n *News) GetAllNews(ctx context.Context) ([]database.News, error) {
	newsList, err := n.q.GetAllNews(ctx)
	if err != nil {
		return nil, err
	}
	return newsList, nil
}

func (n *News) InsertNewsBulk(ctx context.Context, news []database.CreateNewsParams) error {
	tx, err := n.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := n.q.WithTx(tx)

	for _, item := range news {
		if err := qtx.CreateNews(ctx, item); err != nil {
			return err
		}
	}
	return tx.Commit()
}
