package usecase

import (
	"car-sell-buy-system/internal/ads-service/entity"
	"car-sell-buy-system/internal/ads-service/repository"
	"car-sell-buy-system/internal/ads-service/repository/webapi"
	"context"
	"math/big"
)

type AdUseCase struct {
	adRepo repository.AdRepo
	nftApi webapi.NftEthereumWebAPI
}

func NewAdUseCase(adRepo repository.AdRepo) *AdUseCase {
	return &AdUseCase{
		adRepo: adRepo,
	}
}

func (uc *AdUseCase) GetById(ctx context.Context, id int64) (entity.Ad, error) {
	ad, err := uc.adRepo.GetById(ctx, id)
	if err != nil {
		return entity.Ad{}, err
	}

	return ad, nil
}

func (uc *AdUseCase) Store(ctx context.Context, dto entity.AdStoreDTO) (entity.Ad, error) {
	storedAd, err := uc.adRepo.Store(ctx, dto)
	if err != nil {
		return entity.Ad{}, err
	}

	return storedAd, nil
}

func (uc *AdUseCase) List(ctx context.Context, dto entity.AdListDTO) ([]entity.Ad, uint64, error) {
	ads, count, err := uc.adRepo.List(ctx, dto)
	if err != nil {
		return nil, 0, err
	}

	return ads, count, nil
}

func (uc *AdUseCase) HandleFavorite(ctx context.Context, adId, userId int64) error {
	err := uc.adRepo.HandleFavorite(ctx, adId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (uc *AdUseCase) GetTokenInfo(ctx context.Context, tokenId int64) (webapi.NftInfo, error) {
	nftInfo, err := uc.nftApi.GetNftInfo(ctx, big.NewInt(tokenId))
	if err != nil {
		return webapi.NftInfo{}, err
	}

	return nftInfo, nil
}
