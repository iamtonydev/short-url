package url_shortener

import (
	"context"
	"crypto/md5"
	"strings"

	"github.com/jxskiss/base62"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const prefix = "http://localhost:8000/"

func (s *Service) AddShortUrl(ctx context.Context, url string) (string, error) {
	hashBytes := md5.Sum([]byte(url))
	// base62 only contains ASCII symbols so bytes count and characters count will be the same
	hashUrl := base62.EncodeToString(hashBytes[:])

	for i := 0; i < 11; i++ {
		shortUrl := hashUrl[i : i+10]
		err := s.urlsRepository.AddShortUrl(ctx, url, shortUrl)

		if err != nil {
			if s.urlsRepository.IsShortUrlDuplicateError(err) {
				continue
			}

			if s.urlsRepository.IsUrlDuplicateError(err) {
				return "", status.Errorf(codes.InvalidArgument, "this url is already shortened: %v", err.Error())
			}

			return "", err
		}

		var builder strings.Builder
		builder.WriteString(prefix)
		builder.WriteString(shortUrl)
		res := builder.String()

		return res, nil
	}

	return "", status.Errorf(codes.InvalidArgument, "can't shorten this url")
}
