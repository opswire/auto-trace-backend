FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /ads ./cmd/ads

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /ads /ads

EXPOSE 8989

CMD ["/ads"]
