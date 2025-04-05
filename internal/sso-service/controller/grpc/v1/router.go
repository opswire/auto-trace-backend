package v1

import (
	"car-sell-buy-system/internal/sso-service/usecase"
	"car-sell-buy-system/pkg/grpc/api/sso_server_v1"
	"car-sell-buy-system/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouter(grpcServer *grpc.Server, l logger.Interface, u usecase.User) {
	reflection.Register(grpcServer)
	sso_server_v1.RegisterSsoV1Server(grpcServer, &userServer{
		l:  l,
		uc: u,
	})
}
