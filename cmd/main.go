package main

import (
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	authUseCase "live-coding/internal/app/auth"
	userUseCase "live-coding/internal/app/user"
	"live-coding/internal/infrastructure/repository"
	authGrpc "live-coding/internal/interfaces/grpc"
	userGrpc "live-coding/internal/interfaces/grpc"
	"live-coding/proto/auth"
	"live-coding/proto/user"
	"log"
	"net"
)

func main() {

	port := "50051"
	jwtSecret := "mysecret"

	//db connection
	dsn := "host=localhost user=postgres password=admin dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	//repo
	userRepo := repository.NewUserRepository(db)
	roleRightRepo := repository.NewRoleRightsRepository(db)

	//app
	authApp := authUseCase.NewAuthUsecase(userRepo, jwtSecret)
	userApp := userUseCase.NewUserUsecase(userRepo, roleRightRepo, jwtSecret)

	//handler
	authHandler := authGrpc.NewAuthHandler(authApp)
	userHandler := userGrpc.NewUserHandler(userApp)

	//start app
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	auth.RegisterAuthServiceServer(server, authHandler)
	user.RegisterUserServiceServer(server, userHandler)

	log.Printf("server listening at %v", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
