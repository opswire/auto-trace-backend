package webapi

import (
	"car-sell-buy-system/pkg/blockchain/conctracts/carhistory"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"math/big"
	"net/http"
)

const (
	rpcURL        = "https://eth-holesky.g.alchemy.com/v2/CxNXbn3cvcvBDrDOS3lNqK07mHlpeV7Y"
	contractAddr  = "0x67822D3F49d5B59ADd2B5991dd01d3CaaB9A2eB6"
	privateKey    = "0xbd67111e535efece52840724c11172321417fd52c08d3bf39e6d745b492df723"
	chainId       = 17000
	chainName     = "Ethereum holešky"
	chainBasicUrl = "https://holesky.etherscan.io"
)

type NftEthereumWebAPI struct {
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

type NftInfo struct {
	ContractAddr  string        `json:"contract_addr"`
	ChainId       int           `json:"chain_id"`
	ChainName     string        `json:"chain_name"`
	TokenMetadata TokenMetadata `json:"token_metadata"`
	TokenId       int           `json:"token_id"`
	TokenUrl      string        `json:"token_url"`
}

func NewNftEthereumWebAPI() *NftEthereumWebAPI {
	return &NftEthereumWebAPI{}
}

func (*NftEthereumWebAPI) GetNftInfo(ctx context.Context, tokenId *big.Int) (NftInfo, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return NftInfo{}, fmt.Errorf("ошибка подключения к Ethereum: %v", err)
	}

	contract, err := carhistory.NewCarhistory(common.HexToAddress(contractAddr), client)
	if err != nil {
		return NftInfo{}, fmt.Errorf("ошибка создания экземпляра контракта: %v", err)
	}

	uri, err := contract.TokenURI(nil, tokenId)
	if err != nil {
		return NftInfo{}, fmt.Errorf("ошибка получения URI токена: %v", err)
	}

	resp, err := http.Get(uri)
	if err != nil {
		return NftInfo{}, fmt.Errorf("ошибка загрузки метаданных: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return NftInfo{}, fmt.Errorf("ошибка чтения метаданных: %v", err)
	}

	var metadata TokenMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
		return NftInfo{}, fmt.Errorf("ошибка разбора метаданных: %v", err)
	}

	return NftInfo{
		ContractAddr:  contractAddr,
		ChainId:       chainId,
		ChainName:     chainName,
		TokenMetadata: metadata,
		TokenId:       int(tokenId.Int64()),
		TokenUrl:      fmt.Sprintf("%s/nft/%s/%s", chainBasicUrl, contractAddr, tokenId),
	}, nil
}
