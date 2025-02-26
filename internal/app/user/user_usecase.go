package user

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"live-coding/internal/domain"
	"live-coding/internal/infrastructure/repository"
	"strings"
)

type UserUsecase interface {
	ValidateToken(tokenString string) bool
	GetAll(userId uint) ([]domain.UserResponse, error)
	Create(request domain.User, roleId int64) error
	Update(userId uint, name, email, password string) error
	Delete(requesterId, userId uint) error
}

type userUseCase struct {
	userRepo       repository.UserRepository
	roleRightsRepo repository.RoleRightsRepository
	jwtSecret      string
}

func NewUserUsecase(userRepo repository.UserRepository, roleRightsRepo repository.RoleRightsRepository, jwtSecret string) UserUsecase {
	return &userUseCase{
		userRepo:       userRepo,
		roleRightsRepo: roleRightsRepo,
		jwtSecret:      jwtSecret,
	}

}

func (uc *userUseCase) ValidateToken(tokenString string) bool {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return false
	}
	return true
}

func (uc *userUseCase) GetAll(userId uint) ([]domain.UserResponse, error) {

	user, err := uc.userRepo.FindById(userId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	checkPermission, err := uc.roleRightsRepo.CheckPermission(user.RoleId, "users/user", "r_read")
	if err != nil {
		return nil, errors.New("failed to check permission")
	}
	if !checkPermission {
		return nil, errors.New("permission denied")
	}

	users, err := uc.userRepo.GetAll()
	if err != nil {
		return nil, errors.New("failed to get all users")
	}
	return users, nil
}

func (uc *userUseCase) Create(request domain.User, roleId int64) error {
	checkPermission, err := uc.roleRightsRepo.CheckPermission(roleId, "users/user", "r_create")
	if err != nil {
		return errors.New("failed to check permission")
	}
	if !checkPermission {
		return errors.New("permission denied")
	}

	return uc.userRepo.Create(request)
}

func (uc *userUseCase) Update(userId uint, name, email, password string) error {

	user, err := uc.userRepo.FindById(userId)
	if err != nil {
		return errors.New("user not found")
	}

	if email != "" {
		existingUser, err := uc.userRepo.FindByEmail(email)
		if err == nil && existingUser.ID != user.ID {
			return errors.New("email already taken")
		}
	}

	checkPermission, err := uc.roleRightsRepo.CheckPermission(user.RoleId, "users/user", "r_update")
	if err != nil || !checkPermission {
		return errors.New("not allowed to update user")
	}

	return uc.userRepo.Update(userId, name, email, password)

}

func (uc *userUseCase) Delete(requesterId, userId uint) error {
	user, err := uc.userRepo.FindById(requesterId)
	if err != nil {
		return errors.New("user not found")
	}

	checkPermission, err := uc.roleRightsRepo.CheckPermission(user.RoleId, "users/user", "r_delete")
	if err != nil || !checkPermission {
		return errors.New("not allowed to delete user")
	}
	return uc.userRepo.Delete(userId)
}
