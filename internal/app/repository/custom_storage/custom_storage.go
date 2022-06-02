package custom_storage

import (
	"context"
	"errors"

	"github.com/iamtonydev/url-shortener/internal/app/repository"
)

type customStorage struct {
	items map[string]string
}

func NewCustomStorage() repository.IUrlShortenerRepository {
	return &customStorage{
		items: make(map[string]string),
	}
}

func (s *customStorage) AddShortUrl(ctx context.Context, url string, shortUrl string) error {
	if _, found := s.items[shortUrl]; found {
		return nil
	}

	s.items[shortUrl] = url
	return nil
}

func (s *customStorage) GetUrl(ctx context.Context, shortUrl string) (string, error) {
	if val, found := s.items[shortUrl]; found {
		return val, nil
	}

	return "", errors.New("short url not found")
}
