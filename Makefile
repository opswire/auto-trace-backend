include .env
export

.PHONY: migrate-create migrate-up migrate-down

swag-init:
	# ./bin/swag init -g ./internal/ads-service/controller/http/v1/router.go
	 ./bin/swag init -g ./cmd/ads/main.go

swag-fmt:
	 ./bin/swag fmt -g ./cmd/ads/main.go

migrate-create:
	./bin/migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	./bin/migrate -path migrations -database '$(PG_EXPOSE_URL)?sslmode=disable' up

migrate-down:
	./bin/migrate -path migrations -database '$(PG_EXPOSE_URL)?sslmode=disable' down

seed:
	go run ./cmd/ads/seed/seeder.go

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

