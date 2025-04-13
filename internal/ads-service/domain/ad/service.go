package ad

import (
	"car-sell-buy-system/internal/ads-service/domain/nft"
	"context"
	"math/big"
)

type NftRepository interface {
	GetNftInfo(ctx context.Context, tokenId *big.Int) (nft.NFT, error)
}

type Repository interface {
	GetById(ctx context.Context, id int64) (Ad, error)
	Store(ctx context.Context, dto StoreDTO) (Ad, error)
	List(ctx context.Context, dto ListDTO) ([]Ad, uint64, error)
	HandleFavorite(ctx context.Context, adId, userId int64) error
}

type Service struct {
	repository    Repository
	nftRepository NftRepository
}

func NewService(repository Repository, nftRepository NftRepository) *Service {
	return &Service{
		repository:    repository,
		nftRepository: nftRepository,
	}
}

func (s *Service) GetById(ctx context.Context, id int64) (Ad, error) {
	ad, err := s.repository.GetById(ctx, id)
	if err != nil {
		return Ad{}, err
	}

	return ad, nil
}

func (s *Service) Store(ctx context.Context, dto StoreDTO) (Ad, error) {
	storedAd, err := s.repository.Store(ctx, dto)
	if err != nil {
		return Ad{}, err
	}

	return storedAd, nil
}

func (s *Service) List(ctx context.Context, dto ListDTO) ([]Ad, uint64, error) {
	ads, count, err := s.repository.List(ctx, dto)
	if err != nil {
		return nil, 0, err
	}

	return ads, count, nil
}

func (s *Service) HandleFavorite(ctx context.Context, adId, userId int64) error {
	err := s.repository.HandleFavorite(ctx, adId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetTokenInfo(ctx context.Context, tokenId int64) (nft.NFT, error) {
	nftInfo, err := s.nftRepository.GetNftInfo(ctx, big.NewInt(tokenId))
	if err != nil {
		return nft.NFT{}, err
	}

	return nftInfo, nil
}
