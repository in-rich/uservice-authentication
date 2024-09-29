package handlers_test

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	authentication_pb "github.com/in-rich/proto/proto-go/authentication"
	"github.com/in-rich/uservice-authentication/pkg/handlers"
	"github.com/in-rich/uservice-authentication/pkg/models"
	servicesmocks "github.com/in-rich/uservice-authentication/pkg/services/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestListUsers(t *testing.T) {
	testData := []struct {
		name string

		in *authentication_pb.ListUsersRequest

		serviceResponse []*models.User
		serviceErr      error

		expect     *authentication_pb.ListUsersResponse
		expectCode codes.Code
	}{
		{
			name: "ListUsers",
			in: &authentication_pb.ListUsersRequest{
				FirebaseUids: []string{"firebase-uid-1", "firebase-uid-2"},
			},
			serviceResponse: []*models.User{
				{
					PublicIdentifier: "public-identifier-1",
					FirebaseUID:      "firebase-uid-1",
					Email:            "user1@gmail.com",
				},
				{
					PublicIdentifier: "public-identifier-2",
					FirebaseUID:      "firebase-uid-2",
					Email:            "user2@gmail.com",
				},
			},
			expect: &authentication_pb.ListUsersResponse{
				Users: []*authentication_pb.User{
					{
						PublicIdentifier: "public-identifier-1",
						FirebaseUid:      "firebase-uid-1",
						Email:            "user1@gmail.com",
					},
					{
						PublicIdentifier: "public-identifier-2",
						FirebaseUid:      "firebase-uid-2",
						Email:            "user2@gmail.com",
					},
				},
			},
		},
		{
			name: "InternalError",
			in: &authentication_pb.ListUsersRequest{
				FirebaseUids: []string{"firebase-uid-3"},
			},
			serviceErr: errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockListUsersService(t)

			service.On("Exec", context.TODO(), tt.in.FirebaseUids).Return(tt.serviceResponse, tt.serviceErr)

			handler := handlers.NewListUsersHandler(service, monitor.NewDummyGRPCLogger())

			resp, err := handler.ListUsers(context.TODO(), tt.in)

			RequireGRPCCodesEqual(t, err, tt.expectCode)
			require.Equal(t, tt.expect, resp)

			service.AssertExpectations(t)
		})
	}
}
