package handlers

import (
	"context"
	"errors"
	authentication_pb "github.com/in-rich/proto/proto-go/authentication"
	"github.com/in-rich/uservice-authentication/pkg/models"
	"github.com/in-rich/uservice-authentication/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateUserHandler struct {
	authentication_pb.UpdateUserServer
	service services.UpdateUserService
}

func (h *UpdateUserHandler) UpdateUser(ctx context.Context, in *authentication_pb.UpdateUserRequest) (*authentication_pb.User, error) {
	user, err := h.service.Exec(ctx, in.GetToken(), &models.UpdateUser{
		PublicIdentifier: in.GetPublicIdentifier(),
	})

	if err != nil {
		if errors.Is(err, services.ErrUnauthenticated) {
			return nil, status.Errorf(codes.Unauthenticated, "failed to authenticate user: %v", err)
		}
		if errors.Is(err, services.ErrVerifyToken) {
			return nil, status.Errorf(codes.Unauthenticated, "failed to authenticate user: %v", err)
		}
		if errors.Is(err, services.ErrEmailNotVerified) {
			return nil, status.Errorf(codes.PermissionDenied, "failed to authenticate user: %v", err)
		}
		if errors.Is(err, services.ErrInvalidUpdateUser) {
			return nil, status.Errorf(codes.InvalidArgument, "failed to update user: %v", err)
		}
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "failed to update user: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &authentication_pb.User{
		PublicIdentifier: user.PublicIdentifier,
		FirebaseUid:      user.FirebaseUID,
		Email:            user.Email,
	}, nil
}

func NewUpdateUserHandler(service services.UpdateUserService) *UpdateUserHandler {
	return &UpdateUserHandler{
		service: service,
	}
}
