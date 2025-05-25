package ad

import (
	"car-sell-buy-system/internal/ads-service/domain/nft"
	"car-sell-buy-system/pkg/storage/local"
	"context"
	"math/big"
)

type Storage interface {
	Save(file *local.UploadedFile) (string, error)
}

type NftRepository interface {
	GetNftInfo(ctx context.Context, tokenId *big.Int) (nft.NFT, error)
}

type Repository interface {
	GetById(ctx context.Context, id int64) (Ad, error)
	Store(ctx context.Context, dto StoreDTO) (Ad, error)
	Update(ctx context.Context, id int64, dto UpdateDTO) error
	List(ctx context.Context, dto ListDTO) ([]Ad, uint64, error)
	Delete(ctx context.Context, id int64) error
	HandleFavorite(ctx context.Context, adId, userId int64) error
}

type Service struct {
	repository    Repository
	nftRepository NftRepository
	storage       Storage
}

func NewService(
	repository Repository,
	nftRepository NftRepository,
	storage Storage,
) *Service {
	return &Service{
		repository:    repository,
		nftRepository: nftRepository,
		storage:       storage,
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
	if dto.Image != nil {
		pth, err := s.storage.Save(dto.Image)
		if err != nil {
			return Ad{}, err
		}

		dto.CurrentImageUrl = pth
	}

	storedAd, err := s.repository.Store(ctx, dto)
	if err != nil {
		return Ad{}, err
	}

	return storedAd, nil
}

func (s *Service) Update(ctx context.Context, id int64, dto UpdateDTO) error {
	if dto.Image != nil {
		path, err := s.storage.Save(dto.Image)
		if err != nil {
			return err
		}
		dto.CurrentImageUrl = path
	}

	if err := s.repository.Update(ctx, id, dto); err != nil {
		return err
	}

	return nil
}

func (s *Service) List(ctx context.Context, dto ListDTO) ([]Ad, uint64, error) {
	ads, count, err := s.repository.List(ctx, dto)
	if err != nil {
		return nil, 0, err
	}

	return ads, count, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
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
