package main

import (
	"fmt"
	"github.com/in-rich/lib-go/deploy"
	"github.com/in-rich/lib-go/monitor"
	authentication_pb "github.com/in-rich/proto/proto-go/authentication"
	"github.com/in-rich/uservice-authentication/config"
	"github.com/in-rich/uservice-authentication/migrations"
	"github.com/in-rich/uservice-authentication/pkg/dao"
	"github.com/in-rich/uservice-authentication/pkg/handlers"
	"github.com/in-rich/uservice-authentication/pkg/services"
	"github.com/rs/zerolog"
	"os"
)

func getLogger() monitor.GRPCLogger {
	if deploy.IsReleaseEnv() {
		return monitor.NewGCPGRPCLogger(zerolog.New(os.Stdout), "uservice-authentication")
	}

	return monitor.NewConsoleGRPCLogger()
}

func main() {
	logger := getLogger()

	logger.Info("Starting server")
	db, closeDB, err := deploy.OpenDB(config.App.Postgres.DSN)
	if err != nil {
		logger.Fatal(err, "failed to connect to database")
	}
	defer closeDB()

	logger.Info("Running migrations")
	if err := migrations.Migrate(db); err != nil {
		logger.Fatal(err, "failed to migrate")
	}

	depCheck := deploy.DepsCheck{
		Dependencies: func() map[string]error {
			return map[string]error{
				"Postgres": db.Ping(),
			}
		},
		Services: deploy.DepCheckServices{
			"Authenticated": {"Postgres"},
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

	authenticateHandler := handlers.NewAuthenticateHandler(authenticateService, logger)
	getUserHandler := handlers.NewGetUserHandler(getUserService, logger)
	listUsersHandler := handlers.NewListUsersHandler(listUsersService, logger)
	updateUserHandler := handlers.NewUpdateUserHandler(updateUserService, logger)

	logger.Info(fmt.Sprintf("Starting to listen on port %v", config.App.Server.Port))
	listener, server, health := deploy.StartGRPCServer(logger, config.App.Server.Port, depCheck)
	defer deploy.CloseGRPCServer(listener, server)
	go health()

	authentication_pb.RegisterAuthenticateServer(server, authenticateHandler)
	authentication_pb.RegisterGetUserServer(server, getUserHandler)
	authentication_pb.RegisterListUsersServer(server, listUsersHandler)
	authentication_pb.RegisterUpdateUserServer(server, updateUserHandler)

	logger.Info("Server started")
	if err := server.Serve(listener); err != nil {
		logger.Fatal(err, "failed to serve")
	}
}
