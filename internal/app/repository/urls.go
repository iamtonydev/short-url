package repository

//go:generate mockgen --build_flags=--mod=mod -destination=mocks/mock_urls_repository.go -package=mocks . IUrlsRepository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/iamtonydev/url-shortener/internal/app/repository/table"
	"github.com/jmoiron/sqlx"
)

const (
	notFoundErrorMessage          = "sql: Rows are closed"
	urlDuplicateErrorMessage      = "ERROR: duplicate key value violates unique constraint \"urls_url_key\" (SQLSTATE 23505)"
	shortUrlDuplicateErrorMessage = "ERROR: duplicate key value violates unique constraint \"urls_short_url_key\" (SQLSTATE 23505)"
)

type IUrlsRepository interface {
	AddShortUrl(ctx context.Context, url string, shortUrl string) error
	GetUrl(ctx context.Context, shortUrl string) (string, error)
	IsShortUrlDuplicateError(err error) bool
	IsUrlDuplicateError(err error) bool
	IsNotFoundError(err error) bool
}

type urlsRepository struct {
	db *sqlx.DB
}

func NewUrlsRepository(db *sqlx.DB) IUrlsRepository {
	return &urlsRepository{
		db: db,
	}
}

func (r *urlsRepository) AddShortUrl(ctx context.Context, url string, shortUrl string) error {
	builder := sq.Insert(table.UrlsTable).
		PlaceholderFormat(sq.Dollar).
		Columns("url, short_url").
		Values(url, shortUrl)

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

	return row.Err()
}

func (r *urlsRepository) GetUrl(ctx context.Context, shortUrl string) (string, error) {
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

func (r *urlsRepository) IsShortUrlDuplicateError(err error) bool {
	return err.Error() == shortUrlDuplicateErrorMessage
}

func (r *urlsRepository) IsUrlDuplicateError(err error) bool {
	return err.Error() == urlDuplicateErrorMessage
}

func (r *urlsRepository) IsNotFoundError(err error) bool {
	return err.Error() == notFoundErrorMessage
}
