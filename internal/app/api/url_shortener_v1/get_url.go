package url_shortener_v1

import (
	"context"

	desc "github.com/iamtonydev/url-shortener/pkg/url_shortener_v1"
)

func (i *Implementation) GetUrl(ctx context.Context, req *desc.GetUrlRequest) (*desc.GetUrlResponse, error) {
	url, err := i.urlShortenerService.GetUrl(ctx, req.GetShortUrl())
	if err != nil {
		return nil, err
	}

	return &desc.GetUrlResponse{
		Result: &desc.GetUrlResponse_Result{
			Url: url,
		},
	}, nil
}
