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

type AuthenticateHandler struct {
	authentication_pb.AuthenticateServer
	service services.AuthenticateService
	logger  monitor.GRPCLogger
}

func (h *AuthenticateHandler) authenticate(ctx context.Context, in *authentication_pb.AuthenticateRequest) (*authentication_pb.User, error) {
	user, err := h.service.Exec(ctx, in.Token)
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

		return nil, status.Errorf(codes.Internal, "failed to authenticate user: %v", err)
	}

	return &authentication_pb.User{
		PublicIdentifier: user.PublicIdentifier,
		FirebaseUid:      user.FirebaseUID,
		Email:            user.Email,
	}, nil
}

func (h *AuthenticateHandler) Authenticate(ctx context.Context, in *authentication_pb.AuthenticateRequest) (*authentication_pb.User, error) {
	res, err := h.authenticate(ctx, in)
	h.logger.Report(ctx, "Authenticate", err)
	return res, err
}

func NewAuthenticateHandler(service services.AuthenticateService, logger monitor.GRPCLogger) *AuthenticateHandler {
	return &AuthenticateHandler{
		service: service,
		logger:  logger,
	}
}
