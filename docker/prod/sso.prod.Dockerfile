FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /sso ./cmd/sso

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /sso /sso

EXPOSE 8787

CMD ["/sso"]
