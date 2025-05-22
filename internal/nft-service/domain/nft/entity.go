package nft

import (
	"time"
)

type Nft struct {
	TokenId     int64
	Vin         string
	MetadataUrl string
	IsMinted    bool
	CreatedAt   time.Time
	TokenData   TokenData
}

type TokenData struct {
	ContractAddr  string        `json:"contract_addr"`
	ChainId       int           `json:"chain_id"`
	ChainName     string        `json:"chain_name"`
	TokenMetadata TokenMetadata `json:"token_metadata"`
	TokenId       int           `json:"token_id"`
	TokenUrl      string        `json:"token_url"`
	Tx            string        `json:"tx"`
	Records       []Record
}

type Record struct {
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	Company     string    `json:"company"`
	Signature   string    `json:"signature"`
}

type TokenMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Attributes  []struct {
		TraitType string `json:"trait_type"`
		Value     string `json:"value"`
	} `json:"attributes"`
}
