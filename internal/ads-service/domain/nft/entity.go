package nft

type NFT struct {
	ContractAddr  string        `json:"contract_addr"`
	ChainId       int           `json:"chain_id"`
	ChainName     string        `json:"chain_name"`
	TokenMetadata TokenMetadata `json:"token_metadata"`
	TokenId       int           `json:"token_id"`
	TokenUrl      string        `json:"token_url"`
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
