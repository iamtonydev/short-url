package url_shortener

import "github.com/iamtonydev/url-shortener/internal/app/repository"

type Service struct {
	urlShortenerRepository repository.IUrlShortenerRepository
}

func NewUrlShortenerService(urlShortenerRepository repository.IUrlShortenerRepository) *Service {
	return &Service{
		urlShortenerRepository: urlShortenerRepository,
	}
}
