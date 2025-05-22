package nft

import (
	"car-sell-buy-system/internal/nft-service/domain/nft"
	"time"
)

type Response struct {
	TokenId     int64             `json:"token_id"`
	Vin         string            `json:"vin"`
	MetadataUrl string            `json:"metadata_url"`
	IsMinted    bool              `json:"is_minted"`
	CreatedAt   time.Time         `json:"created_at"`
	TokenData   TokenDataResponse `json:"token_data"`
}

type TokenDataResponse struct {
	ContractAddr string `json:"contract_addr"`
	ChainId      int    `json:"chain_id"`
	ChainName    string `json:"chain_name"`
	//TokenMetadata TokenMetadataResponse `json:"token_metadata"`
	TokenId  int          `json:"token_id"`
	TokenUrl string       `json:"token_url"`
	Tx       string       `json:"tx"`
	Records  []nft.Record `json:"records"`
}

type TokenMetadataResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Attributes  []struct {
		TraitType string `json:"trait_type"`
		Value     string `json:"value"`
	} `json:"attributes"`
}

func newResponse(nft nft.Nft) Response {
	return Response{
		TokenId:     nft.TokenId,
		Vin:         nft.Vin,
		MetadataUrl: nft.MetadataUrl,
		IsMinted:    nft.IsMinted,
		CreatedAt:   nft.CreatedAt,
		TokenData: TokenDataResponse{
			ContractAddr: nft.TokenData.ContractAddr,
			ChainId:      nft.TokenData.ChainId,
			ChainName:    nft.TokenData.ChainName,
			//TokenMetadata: TokenMetadataResponse{
			//	Name:        nft.TokenData.TokenMetadata.Name,
			//	Description: nft.TokenData.TokenMetadata.Description,
			//	Image:       nft.TokenData.TokenMetadata.Image,
			//	Attributes:  nft.TokenData.TokenMetadata.Attributes,
			//},
			TokenUrl: nft.TokenData.TokenUrl,
			Tx:       nft.TokenData.Tx,
			Records:  nft.TokenData.Records,
		},
	}
}
