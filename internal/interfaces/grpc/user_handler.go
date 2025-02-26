package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	userApp "live-coding/internal/app/user"
	"live-coding/internal/domain"
	"live-coding/internal/utils"
	"live-coding/proto/user"
	"strconv"
)

type UserHandler struct {
	user.UnimplementedUserServiceServer
	userUseCase userApp.UserUsecase
}

func NewUserHandler(userUseCase userApp.UserUsecase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) validateRequest(ctx context.Context) error {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.InvalidArgument, "missing metadata")
	}

	serviceValues := md.Get("x-link-service")
	if len(serviceValues) == 0 || serviceValues[0] != "be" {
		return status.Error(codes.InvalidArgument, "invalid x-link-service")
	}

	authValues := md.Get("authorization")
	if len(authValues) == 0 {
		return status.Error(codes.InvalidArgument, "invalid authorization")
	}
	token := authValues[0]
	if !h.userUseCase.ValidateToken(token) {
		return status.Error(codes.InvalidArgument, "invalid token")
	}

	return nil
}

func (h *UserHandler) GetAllUsers(ctx context.Context, req *user.GetAllUserRequest) (*user.GetAllUserResponse, error) {
	if err := h.validateRequest(ctx); err != nil {
		return nil, err
	}

	userId, err := utils.ExtractorIdFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "missing user id")
	}

	users, err := h.userUseCase.GetAll(uint(userId))
	if err != nil {
		return &user.GetAllUserResponse{
			Status:  false,
			Message: err.Error(),
		}, nil
	}

	var userResponse []*user.UserResponse
	for _, userRes := range users {
		userResponse = append(userResponse, &user.UserResponse{
			RoleId:     strconv.FormatInt(userRes.RoleId, 10),
			RoleName:   userRes.RoleName,
			Name:       userRes.Name,
			Email:      userRes.Email,
			LastAccess: userRes.LastAccess,
		})
	}

	return &user.GetAllUserResponse{
		Status:  true,
		Message: "Successfully",
		Data:    userResponse,
	}, nil
}

func (h *UserHandler) Create(ctx context.Context, req *user.CreateRequest) (*user.DefaultResponse, error) {

	if err := h.validateRequest(ctx); err != nil {
		return nil, err
	}
	newUser := domain.User{
		RoleId:   req.RoleId,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.userUseCase.Create(newUser, req.RoleId)
	if err != nil {
		return &user.DefaultResponse{
			Status:  false,
			Message: err.Error(),
		}, err
	}
	return &user.DefaultResponse{
		Status:  true,
		Message: "Successfully",
	}, err
}

func (h *UserHandler) Update(ctx context.Context, req *user.UpdateRequest) (*user.DefaultResponse, error) {

	if err := h.validateRequest(ctx); err != nil {
		return nil, err
	}

	userId, err := utils.ExtractorIdFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "missing user id")
	}

	err = h.userUseCase.Update(uint(userId), req.Name, req.Email, req.Password)
	if err != nil {
		return &user.DefaultResponse{
			Status:  false,
			Message: err.Error(),
		}, err
	}
	return &user.DefaultResponse{
		Status:  true,
		Message: "Successfully",
	}, err
}

func (h *UserHandler) Delete(ctx context.Context, req *user.DeleteRequest) (*user.DefaultResponse, error) {
	if err := h.validateRequest(ctx); err != nil {
		return nil, err
	}

	requesterId, err := utils.ExtractorIdFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	userId, err := strconv.ParseUint(req.UserId, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	err = h.userUseCase.Delete(uint(requesterId), uint(userId))
	if err != nil {
		return &user.DefaultResponse{
			Status:  false,
			Message: err.Error(),
		}, err
	}
	return &user.DefaultResponse{
		Status:  true,
		Message: "Successfully",
	}, err
}
