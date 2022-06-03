package url_shortener

import "github.com/iamtonydev/url-shortener/internal/app/repository"

type Service struct {
	urlsRepository repository.IUrlsRepository
}

func NewUrlShortenerService(urlsRepository repository.IUrlsRepository) *Service {
	return &Service{
		urlsRepository: urlsRepository,
	}
}

func NewMockUrlShortenerService(deps ...interface{}) *Service {
	is := Service{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.IUrlsRepository:
			is.urlsRepository = s
		}
	}

	return &is
}
