package grpc

import (
	"context"
	authApp "live-coding/internal/app/auth"
	"live-coding/proto/auth"
)

type AuthHandler struct {
	appAuth authApp.AuthUsecase
	auth.UnimplementedAuthServiceServer
}

func NewAuthHandler(appAuth authApp.AuthUsecase) *AuthHandler {
	return &AuthHandler{appAuth: appAuth}
}

func (h *AuthHandler) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {

	token, err := h.appAuth.Login(request.Email, request.Password)
	if err != nil {
		return &auth.LoginResponse{
			Status:  false,
			Message: "Invalid email or password",
			Data:    &auth.AccessToken{},
		}, nil
	}

	return &auth.LoginResponse{
		Status:  true,
		Message: "Successfully",
		Data: &auth.AccessToken{
			AccessToken: token,
		},
	}, nil
}
