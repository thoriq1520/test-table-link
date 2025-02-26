package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"live-coding/internal/infrastructure/repository"
	"live-coding/internal/utils"
	"time"
)

type AuthUsecase interface {
	Login(email, password string) (string, error)
}

type authUsecase struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthUsecase(repo repository.UserRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{repo: repo, jwtSecret: jwtSecret}
}

func (au *authUsecase) Login(email, password string) (string, error) {

	user, err := au.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if !utils.CheckPassword(user.Password, password) {
		return "", errors.New("wrong password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role_id": user.RoleId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(au.jwtSecret))
	if err != nil {
		return "", err
	}

	err = au.repo.UpdateLastAccess(uint(user.ID), time.Now())
	if err != nil {
		return "", errors.New("failed to update last access")
	}

	return tokenString, nil

}
