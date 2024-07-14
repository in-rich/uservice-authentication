package handlers

import (
	"context"
	authentication_pb "github.com/in-rich/proto/proto-go/authentication"
	"github.com/in-rich/uservice-authentication/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListUsersHandler struct {
	authentication_pb.ListUsersServer
	service services.ListUsersService
}

func (h *ListUsersHandler) ListUsers(ctx context.Context, in *authentication_pb.ListUsersRequest) (*authentication_pb.ListUsersResponse, error) {
	users, err := h.service.Exec(ctx, in.GetFirebaseUids())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	res := &authentication_pb.ListUsersResponse{
		Users: make([]*authentication_pb.User, len(users)),
	}
	for i, user := range users {
		res.Users[i] = &authentication_pb.User{
			PublicIdentifier: user.PublicIdentifier,
			FirebaseUid:      user.FirebaseUID,
			Email:            user.Email,
		}
	}

	return res, nil
}

func NewListUsersHandler(service services.ListUsersService) *ListUsersHandler {
	return &ListUsersHandler{
		service: service,
	}
}
