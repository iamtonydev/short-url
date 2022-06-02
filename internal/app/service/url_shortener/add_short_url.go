package url_shortener

import (
	"context"
	"crypto/md5"

	"github.com/jxskiss/base62"
)

func (s *Service) AddShortUrl(ctx context.Context, url string) (string, error) {
	hashBytes := md5.Sum([]byte(url))
	shortUrl := base62.EncodeToString(hashBytes[:])[:10]
	err := s.urlShortenerRepository.AddShortUrl(ctx, url, shortUrl)
	// TODO: refresh cache case when duplicating
	return shortUrl, err
}
