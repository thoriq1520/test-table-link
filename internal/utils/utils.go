package utils

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
	"strings"
)

func CheckPassword(dbPassword, inputPassword string) bool {
	return dbPassword == inputPassword
}

func ExtractorIdFromContext(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, fmt.Errorf("no metadata")
	}

	authHeader, exist := md["authorization"]
	if !exist {
		return 0, fmt.Errorf("no authorization header")
	}

	tokenParts := strings.Split(authHeader[0], " ")
	if len(tokenParts) != 2 {
		return 0, fmt.Errorf("invalid authorization header")
	}

	userID, err := DecodeJWT(tokenParts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid authorization header")
	}
	return userID, nil
}

func DecodeJWT(tokenString string) (int64, error) {
	jwtSecret := "mysecret"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user id")
	}

	return int64(userId), nil
}
