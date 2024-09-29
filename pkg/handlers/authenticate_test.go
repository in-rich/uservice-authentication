package handlers_test

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	authentication_pb "github.com/in-rich/proto/proto-go/authentication"
	"github.com/in-rich/uservice-authentication/pkg/handlers"
	"github.com/in-rich/uservice-authentication/pkg/models"
	"github.com/in-rich/uservice-authentication/pkg/services"
	servicesmocks "github.com/in-rich/uservice-authentication/pkg/services/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	testData := []struct {
		name string

		in *authentication_pb.AuthenticateRequest

		serviceResponse *models.User
		serviceErr      error

		expect     *authentication_pb.User
		expectCode codes.Code
	}{
		{
			name: "Authenticate",
			in: &authentication_pb.AuthenticateRequest{
				Token: "foo-token",
			},
			serviceResponse: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirebaseUID:      "firebase-uid-1",
				Email:            "user@gmail.com",
			},
			expect: &authentication_pb.User{
				PublicIdentifier: "public-identifier-1",
				FirebaseUid:      "firebase-uid-1",
				Email:            "user@gmail.com",
			},
		},
		{
			name: "Unauthenticated",
			in: &authentication_pb.AuthenticateRequest{
				Token: "foo-token",
			},
			serviceErr: services.ErrUnauthenticated,
			expectCode: codes.Unauthenticated,
		},
		{
			name: "VerifyToken",
			in: &authentication_pb.AuthenticateRequest{
				Token: "foo-token",
			},
			serviceErr: services.ErrVerifyToken,
			expectCode: codes.Unauthenticated,
		},
		{
			name: "EmailNotVerified",
			in: &authentication_pb.AuthenticateRequest{
				Token: "foo-token",
			},
			serviceErr: services.ErrEmailNotVerified,
			expectCode: codes.PermissionDenied,
		},
		{
			name: "InternalError",
			in: &authentication_pb.AuthenticateRequest{
				Token: "foo-token",
			},
			serviceErr: errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockAuthenticateService(t)

			service.On("Exec", context.TODO(), tt.in.Token).Return(tt.serviceResponse, tt.serviceErr)

			handler := handlers.NewAuthenticateHandler(service, monitor.NewDummyGRPCLogger())

			resp, err := handler.Authenticate(context.TODO(), tt.in)

			RequireGRPCCodesEqual(t, err, tt.expectCode)
			require.Equal(t, tt.expect, resp)

			service.AssertExpectations(t)
		})
	}
}
