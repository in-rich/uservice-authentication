package handlers

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	authentication_pb "github.com/in-rich/proto/proto-go/authentication"
	"github.com/in-rich/uservice-authentication/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserHandler struct {
	authentication_pb.GetUserServer
	service services.GetUserService
	logger  monitor.GRPCLogger
}

func (h *GetUserHandler) getUser(ctx context.Context, in *authentication_pb.GetUserRequest) (*authentication_pb.User, error) {
	user, err := h.service.Exec(ctx, in.GetFirebaseUid())
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "failed to get user: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &authentication_pb.User{
		PublicIdentifier: user.PublicIdentifier,
		FirebaseUid:      user.FirebaseUID,
		Email:            user.Email,
	}, nil
}

func (h *GetUserHandler) GetUser(ctx context.Context, in *authentication_pb.GetUserRequest) (*authentication_pb.User, error) {
	res, err := h.getUser(ctx, in)
	h.logger.Report(ctx, "GetUser", err)
	return res, err
}

func NewGetUserHandler(service services.GetUserService, logger monitor.GRPCLogger) *GetUserHandler {
	return &GetUserHandler{
		service: service,
		logger:  logger,
	}
}
