package custom_storage

import (
	"context"
	"errors"
	"sync"

	"github.com/iamtonydev/url-shortener/internal/app/repository"
)

var (
	notFoundError          = errors.New("not found")
	urlDuplicateError      = errors.New("url is already exists")
	shortUrlDuplicateError = errors.New("short url is already exists")
)

type customStorage struct {
	items map[string]string
	mu    sync.Mutex
}

func NewCustomStorage() repository.IUrlsRepository {
	return &customStorage{
		items: make(map[string]string),
		mu:    sync.Mutex{},
	}
}

func (s *customStorage) AddShortUrl(ctx context.Context, url string, shortUrl string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if val, found := s.items[shortUrl]; found {
		if val == url {
			return urlDuplicateError
		}

		return shortUrlDuplicateError
	}

	s.items[shortUrl] = url
	return nil
}

func (s *customStorage) GetUrl(ctx context.Context, shortUrl string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if val, found := s.items[shortUrl]; found {
		return val, nil
	}

	return "", notFoundError
}

func (s *customStorage) IsShortUrlDuplicateError(err error) bool {
	return errors.Is(err, shortUrlDuplicateError)
}

func (s *customStorage) IsUrlDuplicateError(err error) bool {
	return errors.Is(err, urlDuplicateError)
}

func (s *customStorage) IsNotFoundError(err error) bool {
	return errors.Is(err, notFoundError)
}
