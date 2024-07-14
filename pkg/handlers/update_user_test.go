package handlers_test

import (
	"context"
	"errors"
	authentication_pb "github.com/in-rich/proto/proto-go/authentication"
	"github.com/in-rich/uservice-authentication/pkg/handlers"
	"github.com/in-rich/uservice-authentication/pkg/models"
	"github.com/in-rich/uservice-authentication/pkg/services"
	servicesmocks "github.com/in-rich/uservice-authentication/pkg/services/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestUpdateUser(t *testing.T) {
	testData := []struct {
		name string

		in *authentication_pb.UpdateUserRequest

		serviceResponse *models.User
		serviceErr      error

		expect     *authentication_pb.User
		expectCode codes.Code
	}{
		{
			name: "UpdateUser",
			in: &authentication_pb.UpdateUserRequest{
				Token:            "foo-token",
				PublicIdentifier: "public-identifier-2",
			},
			serviceResponse: &models.User{
				PublicIdentifier: "public-identifier-2",
				FirebaseUID:      "firebase-uid-1",
				Email:            "user@gmail.com",
			},
			expect: &authentication_pb.User{
				PublicIdentifier: "public-identifier-2",
				FirebaseUid:      "firebase-uid-1",
				Email:            "user@gmail.com",
			},
		},
		{
			name: "UserNotFound",
			in: &authentication_pb.UpdateUserRequest{
				Token:            "foo-token",
				PublicIdentifier: "public-identifier-2",
			},
			serviceErr: services.ErrUserNotFound,
			expectCode: codes.NotFound,
		},
		{
			name: "EmailNotVerified",
			in: &authentication_pb.UpdateUserRequest{
				Token:            "foo-token",
				PublicIdentifier: "public-identifier-2",
			},
			serviceErr: services.ErrEmailNotVerified,
			expectCode: codes.PermissionDenied,
		},
		{
			name: "VerifyToken",
			in: &authentication_pb.UpdateUserRequest{
				Token:            "foo-token",
				PublicIdentifier: "public-identifier-2",
			},
			serviceErr: services.ErrVerifyToken,
			expectCode: codes.Unauthenticated,
		},
		{
			name: "Unauthenticated",
			in: &authentication_pb.UpdateUserRequest{
				Token:            "foo-token",
				PublicIdentifier: "public-identifier-2",
			},
			serviceErr: services.ErrUnauthenticated,
			expectCode: codes.Unauthenticated,
		},
		{
			name: "InvalidUpdateUser",
			in: &authentication_pb.UpdateUserRequest{
				Token:            "foo-token",
				PublicIdentifier: "public-identifier-2",
			},
			serviceErr: services.ErrInvalidUpdateUser,
			expectCode: codes.InvalidArgument,
		},
		{
			name: "InternalError",
			in: &authentication_pb.UpdateUserRequest{
				Token:            "foo-token",
				PublicIdentifier: "public-identifier-2",
			},
			serviceErr: errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockUpdateUserService(t)

			service.On("Exec", context.TODO(), tt.in.Token, &models.UpdateUser{
				PublicIdentifier: tt.in.PublicIdentifier,
			}).Return(tt.serviceResponse, tt.serviceErr)

			handler := handlers.NewUpdateUserHandler(service)

			resp, err := handler.UpdateUser(context.TODO(), tt.in)

			RequireGRPCCodesEqual(t, err, tt.expectCode)
			require.Equal(t, tt.expect, resp)

			service.AssertExpectations(t)
		})
	}
}
