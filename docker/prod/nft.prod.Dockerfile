FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /nfts ./cmd/nfts

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /nfts /nfts

EXPOSE 8686

CMD ["/nfts"]
