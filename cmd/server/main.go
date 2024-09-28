package main

import (
	"context"
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
	log.Println("Starting server")
	db, closeDB, err := deploy.OpenDB(config.App.Postgres.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer closeDB()

	log.Println("Running migrations")
	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	depCheck := deploy.DepsCheck{
		Dependencies: func() map[string]error {
			_, errAuth := config.AuthClient.GetProjectConfig(context.Background())
			return map[string]error{
				"Postgres":       db.Ping(),
				"Firebase(Auth)": errAuth,
			}
		},
		Services: deploy.DepCheckServices{
			"Authenticated": {"Postgres", "Firebase(Auth)"},
			"GetUser":       {"Postgres"},
			"ListUsers":     {"Postgres"},
			"UpdateUser":    {"Postgres"},
		},
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

	log.Println("Starting to listen on port", config.App.Server.Port)
	listener, server, health := deploy.StartGRPCServer(config.App.Server.Port, depCheck)
	defer deploy.CloseGRPCServer(listener, server)
	go health()

	authentication_pb.RegisterAuthenticateServer(server, authenticateHandler)
	authentication_pb.RegisterGetUserServer(server, getUserHandler)
	authentication_pb.RegisterListUsersServer(server, listUsersHandler)
	authentication_pb.RegisterUpdateUserServer(server, updateUserHandler)

	log.Println("Server started")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
