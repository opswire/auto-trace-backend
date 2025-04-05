package usecase

import (
	"car-sell-buy-system/internal/ads-service/entity"
	"car-sell-buy-system/internal/ads-service/usecase/webapi"
	"context"
)

type (
	Ad interface {
		GetById(ctx context.Context, id int) (*entity.Ad, error)
		Store(ctx context.Context, ad entity.Ad) (entity.Ad, error)
		List(ctx context.Context, dto BasicListRequestDTO) ([]entity.Ad, error)
		HandleFavorite(ctx context.Context, adId, userId int) error
		GetTokenInfo(ctx context.Context, tokenId int) (*webapi.NftInfo, error)
	}

	AdRepo interface {
		GetById(ctx context.Context, id int) (*entity.Ad, error)
		Store(ctx context.Context, ad entity.Ad) (entity.Ad, error)
		List(ctx context.Context, dto BasicListRequestDTO) ([]entity.Ad, error)
		HandleFavorite(ctx context.Context, adId, userId int) error
	}

	CarRepo interface {
		GetById(ctx context.Context, id int) (entity.Car, error)
		GetByVin(ctx context.Context, vin string) (entity.Car, error)
		Store(ctx context.Context, car entity.Car) (entity.Car, error)
	}
)
