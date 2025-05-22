package nft

import (
	"context"
	"math/big"
)

type WebApiRepository interface {
	MintNFT(ctx context.Context, tokenId *big.Int, metadataURI string) (TokenData, error)
	GetNftInfo(ctx context.Context, tokenId *big.Int) (TokenData, error)
	AddServiceRecordToToken(ctx context.Context, tokenId *big.Int, dto AddServiceRecordDTO) (TokenData, error)
	GetServiceRecords(ctx context.Context, tokenId *big.Int) ([]Record, error)
}

type Repository interface {
	StoreNft(ctx context.Context, dto StoreNftDTO) (Nft, error)
	GetNftByVin(ctx context.Context, vin string) (Nft, error)
}

type Service struct {
	repository       Repository
	webApiRepository WebApiRepository
}

func NewService(
	repository Repository,
	webApiRepository WebApiRepository,
) *Service {
	return &Service{
		repository:       repository,
		webApiRepository: webApiRepository,
	}
}

func (s *Service) StoreNft(ctx context.Context, dto StoreNftDTO) (Nft, error) {
	nftInfo, err := s.repository.StoreNft(ctx, dto)
	if err != nil {
		return Nft{}, err
	}

	data, err := s.webApiRepository.MintNFT(ctx, big.NewInt(nftInfo.TokenId), nftInfo.MetadataUrl)
	if err != nil {
		return Nft{}, err
	}

	nftInfo.TokenData = data

	return nftInfo, nil
}

func (s *Service) GetNftByVin(ctx context.Context, vin string) (Nft, error) {
	nftInfo, err := s.repository.GetNftByVin(ctx, vin)
	if err != nil {
		return Nft{}, err
	}

	data, err := s.webApiRepository.GetNftInfo(ctx, big.NewInt(nftInfo.TokenId))
	if err != nil {
		return Nft{}, err
	}

	records, err := s.webApiRepository.GetServiceRecords(ctx, big.NewInt(nftInfo.TokenId))
	if err != nil {
		return Nft{}, err
	}

	nftInfo.TokenData = data
	nftInfo.TokenData.Records = records

	return nftInfo, nil
}

func (s *Service) AddServiceRecordByVin(ctx context.Context, vin string, dto AddServiceRecordDTO) (Nft, error) {
	nftInfo, err := s.repository.GetNftByVin(ctx, vin)
	if err != nil {
		return Nft{}, err
	}

	data, err := s.webApiRepository.AddServiceRecordToToken(ctx, big.NewInt(nftInfo.TokenId), dto)
	if err != nil {
		return Nft{}, err
	}

	nftInfo.TokenData = data

	return nftInfo, nil
}
