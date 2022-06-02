package url_shortener_v1

import (
	"github.com/iamtonydev/url-shortener/internal/app/service/url_shortener"
	desc "github.com/iamtonydev/url-shortener/pkg/url_shortener_v1"
)

type Implementation struct {
	desc.UnimplementedUrlShortenerV1Server

	urlShortenerService *url_shortener.Service
}

func NewUrlShortenerV1(urlShortenerService *url_shortener.Service) *Implementation {
	return &Implementation{
		desc.UnimplementedUrlShortenerV1Server{},

		urlShortenerService,
	}
}
