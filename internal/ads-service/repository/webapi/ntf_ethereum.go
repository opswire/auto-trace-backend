package webapi

import (
	"car-sell-buy-system/internal/ads-service/domain/nft"
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

func NewNftEthereumWebAPI() *NftEthereumWebAPI {
	return &NftEthereumWebAPI{}
}

func (*NftEthereumWebAPI) GetNftInfo(ctx context.Context, tokenId *big.Int) (nft.NFT, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nft.NFT{}, fmt.Errorf("ошибка подключения к Ethereum: %v", err)
	}

	contract, err := carhistory.NewCarhistory(common.HexToAddress(contractAddr), client)
	if err != nil {
		return nft.NFT{}, fmt.Errorf("ошибка создания экземпляра контракта: %v", err)
	}

	uri, err := contract.TokenURI(nil, tokenId)
	if err != nil {
		return nft.NFT{}, fmt.Errorf("ошибка получения URI токена: %v", err)
	}

	resp, err := http.Get(uri)
	if err != nil {
		return nft.NFT{}, fmt.Errorf("ошибка загрузки метаданных: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nft.NFT{}, fmt.Errorf("ошибка чтения метаданных: %v", err)
	}

	var metadata nft.TokenMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
		return nft.NFT{}, fmt.Errorf("ошибка разбора метаданных: %v", err)
	}

	return nft.NFT{
		ContractAddr:  contractAddr,
		ChainId:       chainId,
		ChainName:     chainName,
		TokenMetadata: metadata,
		TokenId:       int(tokenId.Int64()),
		TokenUrl:      fmt.Sprintf("%s/nft/%s/%s", chainBasicUrl, contractAddr, tokenId),
	}, nil
}
