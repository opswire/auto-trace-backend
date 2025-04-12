package usecase

import (
	"car-sell-buy-system/internal/ads-service/entity"
	"car-sell-buy-system/internal/ads-service/repository/webapi"
	"context"
)

type (
	Ad interface {
		GetById(ctx context.Context, id int64) (entity.Ad, error)
		Store(ctx context.Context, dto entity.AdStoreDTO) (entity.Ad, error)
		List(ctx context.Context, dto entity.AdListDTO) ([]entity.Ad, uint64, error)
		HandleFavorite(ctx context.Context, adId, userId int64) error
		GetTokenInfo(ctx context.Context, tokenId int64) (webapi.NftInfo, error)
	}
)
