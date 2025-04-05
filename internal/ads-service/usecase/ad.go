package usecase

import (
	"car-sell-buy-system/internal/ads-service/entity"
	"car-sell-buy-system/internal/ads-service/usecase/webapi"
	"context"
	"math/big"
)

type AdUseCase struct {
	adRepo  AdRepo
	carRepo CarRepo
	nftApi  webapi.NftEthereumWebAPI
}

func NewAdUseCase(adRepo AdRepo, carRepo CarRepo) *AdUseCase {
	return &AdUseCase{
		adRepo:  adRepo,
		carRepo: carRepo,
	}
}

func (uc *AdUseCase) GetById(ctx context.Context, id int) (*entity.Ad, error) {
	ad, err := uc.adRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (uc *AdUseCase) Store(ctx context.Context, ad entity.Ad) (entity.Ad, error) {
	car, err := uc.carRepo.GetByVin(ctx, ad.Car.Vin)
	if err != nil {
		return entity.Ad{}, err
	}

	if car == (entity.Car{}) {
		ad.Car, err = uc.carRepo.Store(ctx, ad.Car)
		if err != nil {
			return entity.Ad{}, err
		}
	} else {
		ad.Car = car
	}

	storedAd, err := uc.adRepo.Store(ctx, ad)
	if err != nil {
		return entity.Ad{}, err
	}

	return storedAd, nil
}

func (uc *AdUseCase) List(ctx context.Context, dto BasicListRequestDTO) ([]entity.Ad, error) {
	ads, err := uc.adRepo.List(ctx, dto)
	if err != nil {
		return nil, err
	}

	return ads, nil
}

func (uc *AdUseCase) HandleFavorite(ctx context.Context, adId, userId int) error {
	err := uc.adRepo.HandleFavorite(ctx, adId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (uc *AdUseCase) GetTokenInfo(ctx context.Context, tokenId int) (*webapi.NftInfo, error) {
	nftInfo, err := uc.nftApi.GetNftInfo(ctx, big.NewInt(int64(tokenId)))
	if err != nil {
		return nil, err
	}

	return nftInfo, nil
}
