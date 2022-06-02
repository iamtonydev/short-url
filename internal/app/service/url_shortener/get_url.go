package url_shortener

import "context"

func (s *Service) GetUrl(ctx context.Context, shortUrl string) (string, error) {
	return s.urlShortenerRepository.GetUrl(ctx, shortUrl)
}
