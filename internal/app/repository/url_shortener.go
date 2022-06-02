package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/iamtonydev/url-shortener/internal/app/repository/table"
	"github.com/jmoiron/sqlx"
)

type IUrlShortenerRepository interface {
	AddShortUrl(ctx context.Context, url string, shortUrl string) error
	GetUrl(ctx context.Context, shortUrl string) (string, error)
}

type urlShortenerRepository struct {
	db *sqlx.DB
}

func NewUrlShortenerRepository(db *sqlx.DB) IUrlShortenerRepository {
	return &urlShortenerRepository{
		db: db,
	}
}

func (r *urlShortenerRepository) AddShortUrl(ctx context.Context, url string, shortUrl string) error {
	builder := sq.Insert(table.UrlsTable).
		PlaceholderFormat(sq.Dollar).
		Columns("url, short_url").
		Values(url, shortUrl).
		Suffix("returning short_url")

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	row, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer row.Close()

	row.Next()
	err = row.Scan()
	if err != nil {
		// TODO: handle unique constraint error for short_url
		return err
	}

	return nil
}

func (r *urlShortenerRepository) GetUrl(ctx context.Context, shortUrl string) (string, error) {
	builder := sq.Select("url").
		PlaceholderFormat(sq.Dollar).
		From(table.UrlsTable).
		Where(sq.Eq{"short_url": shortUrl}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return "", err
	}

	row, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return "", err
	}
	defer row.Close()

	row.Next()
	var url string
	err = row.Scan(&url)
	if err != nil {
		return "", err
	}

	return url, nil
}
