include .env
export

.PHONY: migrate-create migrate-up migrate-down

up:
	docker compose up -d
.PHONY: up

swag-init:
	 swag init -g ./internal/ads-service/controller/http/router.go

swag-fmt:
	 swag fmt -g ./cmd/ads/main.go

nfts-swag-init:
	swag init -g ./internal/nft-service/controller/http/router.go

nfts-swag-fmt:
	swag fmt -g ./cmd/nfts/main.go

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	go run ./cmd/command/root.go migrate-up

migrate-down:
	go run ./cmd/command/root.go migrate-down

seed:
	go run ./cmd/ads/seed/seeder.go

.PHONY: test
test: ### run test
	go test -v -race -covermode atomic -coverprofile=cover.out ./internal/ads-service...

.PHONY: tool-test
tool-test: ### run test
	 go tool cover -html=./cover.out -o ./cover.html

.PHONY: mock-ad
mock-ad: ### run mockgen
	mockgen -source ./internal/ads-service/domain/ad/service.go -package ad_test > ./internal/ads-service/domain/ad/mock_service_test.go

.PHONY: mock-chat
mock-chat: ### run mockgen
	mockgen -source ./internal/ads-service/domain/chat/service.go -package chat_test > ./internal/ads-service/domain/chat/mock_service_test.go

.PHONY: mock-payment
mock-payment: ### run mockgen
	mockgen -source ./internal/ads-service/domain/payment/service.go -package payment_test > ./internal/ads-service/domain/payment/mock_service_test.go

# GRPC

# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN := $(CURDIR)/bin

# Добавляем bin в текущей директории в PATH при запуске protoc
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

# Путь до protobuf файлов
PROTO_PATH := $(CURDIR)/api

# Путь до сгенеренных .pb.go файлов
PKG_PROTO_PATH := $(CURDIR)/pkg/grpc

# устанавливаем необходимые плагины
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	$(info Installing binary dependencies...)

	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/swaggo/swag/cmd/swag@latest

	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# генерация .go файлов с помощью protoc
.protoc-generate:
	mkdir -p $(PKG_PROTO_PATH)
	$(PROTOC) --proto_path=$(CURDIR) \
	--go_out=$(PKG_PROTO_PATH) --go_opt paths=source_relative \
	--go-grpc_out=$(PKG_PROTO_PATH) --go-grpc_opt paths=source_relative \
	$(PROTO_PATH)/sso_server_v1/sso.proto # \
	# $(PROTO_PATH)/notes/messages.proto

# go mod tidy
.tidy:
	GOBIN=$(LOCAL_BIN) go mod tidy

# Генерация кода из protobuf
grpc-api-generate: .bin-deps .protoc-generate .tidy

# Объявляем, что текущие команды не являются файлами и
# интсрументируем Makefile не искать изменения в файловой системе
.PHONY: \
	.bin-deps \
	.protoc-generate \
	.tidy \
	generate

# Blockchain

abi-generate:
	abigen --abi pkg/blockchain/conctracts/carhistory/CarHistoryNFT.abi --pkg carhistory --out ./pkg/blockchain/conctracts/carhistory/CarHistoryNFT.go
.PHONY: abi-generate

