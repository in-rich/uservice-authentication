package main

import (
	"fmt"
	"github.com/in-rich/lib-go/deploy"
	authentication_pb "github.com/in-rich/proto/proto-go/authentication"
	"github.com/in-rich/uservice-authentication/config"
	"github.com/in-rich/uservice-authentication/migrations"
	"github.com/in-rich/uservice-authentication/pkg/dao"
	"github.com/in-rich/uservice-authentication/pkg/handlers"
	"github.com/in-rich/uservice-authentication/pkg/services"
	"log"
)

func main() {
	db, closeDB := deploy.OpenDB(config.App.Postgres.DSN)
	defer closeDB()

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	getUsersDAO := dao.NewGetUserRepository(db)
	listUsersDAO := dao.NewListUsersRepository(db)
	createUserDAO := dao.NewCreateUserRepository(db)
	updateUserDAO := dao.NewUpdateUserRepository(db)

	authenticateService := services.NewAuthenticateService(config.AuthClient, getUsersDAO)
	getUserService := services.NewGetUserService(config.AuthClient, getUsersDAO)
	listUsersService := services.NewListUsersService(config.AuthClient, listUsersDAO)
	updateUserService := services.NewUpdateUserService(authenticateService, createUserDAO, updateUserDAO)

	authenticateHandler := handlers.NewAuthenticateHandler(authenticateService)
	getUserHandler := handlers.NewGetUserHandler(getUserService)
	listUsersHandler := handlers.NewListUsersHandler(listUsersService)
	updateUserHandler := handlers.NewUpdateUserHandler(updateUserService)

	listener, server := deploy.StartGRPCServer(fmt.Sprintf(":%d", config.App.Server.Port), "authentication")
	defer deploy.CloseGRPCServer(listener, server)

	authentication_pb.RegisterAuthenticateServer(server, authenticateHandler)
	authentication_pb.RegisterGetUserServer(server, getUserHandler)
	authentication_pb.RegisterListUsersServer(server, listUsersHandler)
	authentication_pb.RegisterUpdateUserServer(server, updateUserHandler)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
