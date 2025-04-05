package v1

import (
	"car-sell-buy-system/internal/sso-service/usecase"
	"car-sell-buy-system/pkg/auth"
	"car-sell-buy-system/pkg/grpc/api/sso_server_v1"
	"car-sell-buy-system/pkg/logger"
	"context"
	"strconv"
	"strings"
)

type userServer struct {
	sso_server_v1.UnimplementedSsoV1Server
	uc usecase.User
	l  logger.Interface
}

func (s *userServer) VerifyToken(ctx context.Context, req *sso_server_v1.VerifyTokenRequest) (*sso_server_v1.VerifyTokenResponse, error) {
	s.l.Info("token: "+req.GetToken(), req.Token)

	parts := strings.Split(req.GetToken(), " ")
	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return &sso_server_v1.VerifyTokenResponse{
			UserId: 0,
			Valid:  false,
		}, nil
	}

	claims, err := auth.ParseJWT(parts[1])
	if err != nil {
		return &sso_server_v1.VerifyTokenResponse{
			UserId: 0,
			Valid:  false,
		}, nil
	}

	userId, err := strconv.Atoi(claims.ID)
	if err != nil {
		return &sso_server_v1.VerifyTokenResponse{
			UserId: 0,
			Valid:  false,
		}, nil
	}

	return &sso_server_v1.VerifyTokenResponse{
		UserId: int64(userId),
		Valid:  true,
	}, nil
}
