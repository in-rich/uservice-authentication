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

func TestGetUser(t *testing.T) {
	testData := []struct {
		name string

		in *authentication_pb.GetUserRequest

		serviceResponse *models.User
		serviceErr      error

		expect     *authentication_pb.User
		expectCode codes.Code
	}{
		{
			name: "GetUser",
			in: &authentication_pb.GetUserRequest{
				FirebaseUid: "firebase-uid-1",
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
			name: "UserNotFound",
			in: &authentication_pb.GetUserRequest{
				FirebaseUid: "firebase-uid-2",
			},
			serviceErr: services.ErrUserNotFound,
			expectCode: codes.NotFound,
		},
		{
			name: "InternalError",
			in: &authentication_pb.GetUserRequest{
				FirebaseUid: "firebase-uid-3",
			},
			serviceErr: errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockGetUserService(t)

			service.On("Exec", context.TODO(), tt.in.FirebaseUid).Return(tt.serviceResponse, tt.serviceErr)

			handler := handlers.NewGetUserHandler(service, monitor.NewDummyGRPCLogger())

			resp, err := handler.GetUser(context.TODO(), tt.in)

			RequireGRPCCodesEqual(t, err, tt.expectCode)
			require.Equal(t, tt.expect, resp)

			service.AssertExpectations(t)
		})
	}
}
