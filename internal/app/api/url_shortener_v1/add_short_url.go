package url_shortener_v1

import (
	"context"

	desc "github.com/iamtonydev/url-shortener/pkg/url_shortener_v1"
)

func (i *Implementation) AddShortUrl(ctx context.Context, req *desc.AddShortUrlRequest) (*desc.AddShortUrlResponse, error) {
	shortUrl, err := i.urlShortenerService.AddShortUrl(ctx, req.GetUrl())
	if err != nil {
		return nil, err
	}

	return &desc.AddShortUrlResponse{
		Result: &desc.AddShortUrlResponse_Result{
			ShortUrl: shortUrl,
		},
	}, nil
}
