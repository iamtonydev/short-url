package url_shortener_v1

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	urlRepoMocks "github.com/iamtonydev/url-shortener/internal/app/repository/mocks"
	"github.com/iamtonydev/url-shortener/internal/app/service/url_shortener"
	desc "github.com/iamtonydev/url-shortener/pkg/url_shortener_v1"
	"github.com/stretchr/testify/require"
)

func TestAddShortUrl(t *testing.T) {
	var (
		ctx      = context.Background()
		mockCtrl = gomock.NewController(t)

		url = gofakeit.URL()

		validReq = &desc.AddShortUrlRequest{Url: url}

		errRepo = "this url is already shortened"
	)

	urlRepoMock := urlRepoMocks.NewMockIUrlsRepository(mockCtrl)
	gomock.InOrder(
		urlRepoMock.EXPECT().AddShortUrl(ctx, url, gomock.Any()).Return(nil).Times(1),
		urlRepoMock.EXPECT().AddShortUrl(ctx, url, gomock.Any()).Return(errors.New(errRepo)).Times(1),
		urlRepoMock.EXPECT().IsShortUrlDuplicateError(gomock.Any()).Return(false).Times(1),
		urlRepoMock.EXPECT().IsUrlDuplicateError(gomock.Any()).Return(false).Times(1),
	)

	api := newMockUrlShortenerV1(Implementation{
		urlShortenerService: url_shortener.NewMockUrlShortenerService(
			urlRepoMock,
		),
	})

	t.Run("success case", func(t *testing.T) {
		_, err := api.AddShortUrl(ctx, validReq)
		require.Nil(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		_, err := api.AddShortUrl(ctx, validReq)
		require.Error(t, err)
		require.Equal(t, errRepo, err.Error())
	})
}
