package main

import (
	"car-sell-buy-system/pkg/blockchain/conctracts/carhistory"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	rpcURL       = "https://eth-holesky.g.alchemy.com/v2/CxNXbn3cvcvBDrDOS3lNqK07mHlpeV7Y"
	contractAddr = "0x67822D3F49d5B59ADd2B5991dd01d3CaaB9A2eB6"
	privateKey   = "0xbd67111e535efece52840724c11172321417fd52c08d3bf39e6d745b492df723"
	chainId      = 17000
)

func main() {
	fmt.Println(GetNftInfo(big.NewInt(5552)))
	return
	// Подключение к сети
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Загрузка приватного ключа
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey[2:])
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Получение публичного адреса
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Invalid key")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Получение nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// Установка газа и цены газа
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	// Настройка транзакции
	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, big.NewInt(chainId))
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // Нет оплаты в ETH

	// Установите цену газа ниже текущей средней
	gasPrice = gasPrice.Div(gasPrice, big.NewInt(2))
	auth.GasPrice = gasPrice

	// Подключение к контракту
	contract, err := carhistory.NewCarhistory(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatalf("Failed to load contract: %v", err)
	}

	// Параметры вызова mintNFT
	toAddress := common.HexToAddress("0xbe012aE4Ab6BdC891884fb5a7763A5F7dbC26eef")                                    // Получатель токена
	tokenId := big.NewInt(1)                                                                                          // Уникальный ID токена
	tokenURI := "https://ivory-total-catshark-439.mypinata.cloud/ipfs/QmdjzwyKPUv7EhW38ywuJKP25or6TkNdmBGodazoUArQNM" // URI метаданных токена

	tx, err := contract.MintCar(auth, toAddress, tokenId, tokenURI)
	if err != nil {
		log.Fatalf("Failed to mint NFT: %v", err)
	}

	log.Printf("Transaction sent: %s", tx.Hash().Hex())
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

func GetNftInfo(tokenId *big.Int) TokenMetadata {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к Ethereum: %v", err)
	}

	// Получаем URI токена
	contract, err := carhistory.NewCarhistory(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatalf("Ошибка создания экземпляра контракта: %v", err)
	}

	uri, err := contract.TokenURI(nil, tokenId)
	if err != nil {
		log.Fatalf("Ошибка получения URI токена: %v", err)
	}

	// Загружаем метаданные из URI
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalf("Ошибка загрузки метаданных: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка чтения метаданных: %v", err)
	}

	var metadata TokenMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
		log.Fatalf("Ошибка разбора метаданных: %v", err)
	}

	return metadata
}
