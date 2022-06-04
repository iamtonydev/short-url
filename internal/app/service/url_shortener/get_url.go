package url_shortener

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetUrl(ctx context.Context, shortUrl string) (string, error) {
	url, err := s.urlsRepository.GetUrl(ctx, shortUrl)
	if err != nil {
		if s.urlsRepository.IsNotFoundError(err) {
			return "", status.Errorf(codes.NotFound, "short url not found: %v", err.Error())
		}

		return "", err
	}

	return url, err
}
