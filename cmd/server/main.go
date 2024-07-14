package main

import (
	"database/sql"
	"fmt"
	authentication_pb "github.com/in-rich/proto/proto-go/authentication"
	"github.com/in-rich/uservice-authentication/config"
	"github.com/in-rich/uservice-authentication/migrations"
	"github.com/in-rich/uservice-authentication/pkg/dao"
	"github.com/in-rich/uservice-authentication/pkg/handlers"
	"github.com/in-rich/uservice-authentication/pkg/services"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.App.Postgres.DSN)))
	db := bun.NewDB(sqldb, pgdialect.New())

	defer func() {
		_ = db.Close()
		_ = sqldb.Close()
	}()

	err := db.Ping()
	for i := 0; i < 10 && err != nil; i++ {
		time.Sleep(1 * time.Second)
		err = db.Ping()
	}

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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.App.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	defer func() {
		server.GracefulStop()
		_ = listener.Close()
	}()

	authentication_pb.RegisterAuthenticateServer(server, authenticateHandler)
	authentication_pb.RegisterGetUserServer(server, getUserHandler)
	authentication_pb.RegisterListUsersServer(server, listUsersHandler)
	authentication_pb.RegisterUpdateUserServer(server, updateUserHandler)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
