package handlers

import (
	"context"
	"errors"
	authentication_pb "github.com/in-rich/proto/proto-go/authentication"
	"github.com/in-rich/uservice-authentication/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserHandler struct {
	authentication_pb.GetUserServer
	service services.GetUserService
}

func (h *GetUserHandler) GetUser(ctx context.Context, in *authentication_pb.GetUserRequest) (*authentication_pb.User, error) {
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

func NewGetUserHandler(service services.GetUserService) *GetUserHandler {
	return &GetUserHandler{
		service: service,
	}
}
