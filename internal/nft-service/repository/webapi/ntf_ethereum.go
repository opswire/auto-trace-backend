package webapi

import (
	"car-sell-buy-system/internal/nft-service/domain/nft"
	"car-sell-buy-system/pkg/blockchain/conctracts/carhistory"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

const (
	rpcURL         = "https://eth-holesky.g.alchemy.com/v2/CxNXbn3cvcvBDrDOS3lNqK07mHlpeV7Y"
	contractAddr   = "0x631884DC264999f02E0CFf7D36Cd12Dbd7aEae8f"
	privateKeyAddr = "672c7a28eea558990b26fc49ffe7aeda99a7d5f13d2e5056ce288afac8eb00ff"
	chainId        = 17000
	chainName      = "Ethereum holešky"
	chainBasicUrl  = "https://holesky.etherscan.io"
	toAddress      = "0x36b46587441b0CC2De26343233F5DC5539F5D3D9"
)

type NftEthereumWebAPI struct {
}

func NewNftEthereumWebAPI() *NftEthereumWebAPI {
	return &NftEthereumWebAPI{}
}

func (*NftEthereumWebAPI) MintNFT(ctx context.Context, tokenId *big.Int, metadataURI string) (nft.TokenData, error) {
	// 1. Подключение к Ethereum-сети
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nft.TokenData{}, fmt.Errorf("ошибка подключения к Ethereum: %v", err)
	}

	// 2. Создание экземпляра контракта
	contract, err := carhistory.NewCarhistory(common.HexToAddress(contractAddr), client)
	if err != nil {
		return nft.TokenData{}, fmt.Errorf("ошибка создания экземпляра контракта: %v", err)
	}

	// 3. Настройка приватного ключа для подписания транзакций
	privateKey, err := crypto.HexToECDSA(privateKeyAddr)
	if err != nil {
		return nft.TokenData{}, fmt.Errorf("ошибка загрузки приватного ключа: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nft.TokenData{}, fmt.Errorf("некорректный тип публичного ключа")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 4. Получение nonce для отправителя
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nft.TokenData{}, fmt.Errorf("ошибка получения nonce: %v", err)
	}

	// 5. Получение текущей цены газа
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nft.TokenData{}, fmt.Errorf("ошибка получения цены газа: %v", err)
	}

	// 6. Создание авторизатора для транзакции
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	// 7. Вызов функции mintCar
	tx, err := contract.MintCar(auth, common.HexToAddress(toAddress), tokenId, metadataURI)
	if err != nil {
		return nft.TokenData{}, fmt.Errorf("ошибка выпуска токена: %v", err)
	}

	fmt.Printf("Транзакция отправлена! Хэш транзакции: %s\n", tx.Hash().Hex())

	return nft.TokenData{
		ContractAddr:  contractAddr,
		ChainId:       chainId,
		ChainName:     chainName,
		TokenMetadata: nft.TokenMetadata{},
		TokenId:       int(tokenId.Int64()),
		TokenUrl:      fmt.Sprintf("%s/nft/%s/%s", chainBasicUrl, contractAddr, tokenId),
		Tx:            tx.Hash().Hex(),
	}, nil
}

func (*NftEthereumWebAPI) GetNftInfo(ctx context.Context, tokenId *big.Int) (nft.TokenData, error) {
	// for metadata 1
	//client, err := ethclient.Dial(rpcURL)
	//if err != nil {
	//	return nft.TokenData{}, fmt.Errorf("ошибка подключения к Ethereum: %v", err)
	//}
	//
	//contract, err := carhistory.NewCarhistory(common.HexToAddress(contractAddr), client)
	//if err != nil {
	//	return nft.TokenData{}, fmt.Errorf("ошибка создания экземпляра контракта: %v", err)
	//}
	//
	//uri, err := contract.TokenURI(nil, big.NewInt(123))
	//if err != nil {
	//	return nft.TokenData{}, fmt.Errorf("ошибка получения URI токена: %v", err)
	//}
	//
	//fmt.Println("11 uri: ", uri)
	//fmt.Println("2 tokenId: ", tokenId)
	//resp, err := http.Get(uri)
	//if err != nil {
	//	return nft.TokenData{}, fmt.Errorf("ошибка загрузки метаданных: %v", err)
	//}
	//defer resp.Body.Close()
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return nft.TokenData{}, fmt.Errorf("ошибка чтения метаданных: %v", err)
	//}
	//
	//fmt.Println(body)
	//var metadata nft.TokenMetadata
	//if err = json.Unmarshal(body, &metadata); err != nil {
	//	return nft.TokenData{}, fmt.Errorf("ошибка разбора метаданных: %v", err)
	//}

	return nft.TokenData{
		ContractAddr: contractAddr,
		ChainId:      chainId,
		ChainName:    chainName,
		TokenId:      int(tokenId.Int64()),
		TokenUrl:     fmt.Sprintf("%s/nft/%s/%s", chainBasicUrl, contractAddr, tokenId),
	}, nil
}

func (*NftEthereumWebAPI) AddServiceRecordToToken(
	ctx context.Context,
	tokenId *big.Int,
	dto nft.AddServiceRecordDTO,
) (nft.TokenData, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nft.TokenData{}, fmt.Errorf("ошибка подключения к Ethereum: %v", err)
	}

	contract, err := carhistory.NewCarhistory(common.HexToAddress(contractAddr), client)
	if err != nil {
		return nft.TokenData{}, fmt.Errorf("ошибка создания экземпляра контракта: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyAddr)
	if err != nil {
		return nft.TokenData{}, fmt.Errorf("ошибка загрузки приватного ключа: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nft.TokenData{}, fmt.Errorf("некорректный тип публичного ключа")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nft.TokenData{}, fmt.Errorf("ошибка получения nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nft.TokenData{}, fmt.Errorf("ошибка получения цены газа: %v", err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	tx, err := contract.AddServiceRecord(
		auth,
		tokenId,
		dto.Description,
		dto.Company,
		dto.Signature,
	)
	if err != nil {
		return nft.TokenData{}, fmt.Errorf("ошибка добавления записи о ТО: %v", err)
	}

	return nft.TokenData{
		ContractAddr:  contractAddr,
		ChainId:       chainId,
		ChainName:     chainName,
		TokenMetadata: nft.TokenMetadata{},
		TokenId:       int(tokenId.Int64()),
		TokenUrl:      fmt.Sprintf("%s/nft/%s/%s", chainBasicUrl, contractAddr, tokenId),
		Tx:            fmt.Sprintf("Транзакция отправлена! Хэш транзакции: %s\n", tx.Hash().Hex()),
	}, nil
}

func (*NftEthereumWebAPI) GetServiceRecords(ctx context.Context, tokenId *big.Int) ([]nft.Record, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к Ethereum: %v", err)
	}

	contract, err := carhistory.NewCarhistory(common.HexToAddress(contractAddr), client)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания экземпляра контракта: %v", err)
	}

	records, err := contract.GetServiceRecords(nil, tokenId)
	if err != nil {
		return nil, fmt.Errorf("ошибка добавления записи о ТО: %v", err)
	}

	fmt.Printf("История получена! История: %s\n", records)

	nftRecords := make([]nft.Record, 0, len(records))
	for _, record := range records {
		if record.Description != "" {
			nftRecords = append(nftRecords, nft.Record{
				CreatedAt:   time.Unix(record.Timestamp.Int64(), 0).UTC(),
				Description: record.Description,
				Company:     record.Company,
				Signature:   record.Signature,
			})
		}
	}

	return nftRecords, nil
}
