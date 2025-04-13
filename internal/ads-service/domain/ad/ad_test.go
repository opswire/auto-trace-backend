package ad_test

import (
	"car-sell-buy-system/internal/ads-service/domain/ad"
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

var errInternalServErr = errors.New("internal server error")

type test struct {
	name string
	mock func()
	res  interface{}
	err  error
}

func mockService(t *testing.T) (*ad.Service, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repository := NewMockRepository(mockCtl)
	nftRepository := NewMockNftRepository(mockCtl)

	service := ad.NewService(repository, nftRepository)

	return service, repository
}

func TestList(t *testing.T) {
	t.Parallel()

	service, repo := mockService(t)

	dto := ad.ListDTO{}

	fakeAd := ad.Ad{}

	tests := []test{
		{
			name: "successful ad list",
			mock: func() {
				repo.
					EXPECT().
					List(context.Background(), dto).
					Return([]ad.Ad{fakeAd}, uint64(1), nil)
			},
			res: []ad.Ad{fakeAd},
			err: nil,
		},
		{
			name: "empty result",
			mock: func() {
				repo.EXPECT().
					List(context.Background(), dto).
					Return(nil, uint64(0), nil)
			},
			res: []ad.Ad(nil),
			err: nil,
		},
		{
			name: "result with error",
			mock: func() {
				repo.
					EXPECT().
					List(context.Background(), dto).
					Return(nil, uint64(0), errInternalServErr)
			},
			res: []ad.Ad(nil),
			err: errInternalServErr,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			res, _, err := service.List(context.Background(), dto)

			require.Equal(t, res, localTc.res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}
