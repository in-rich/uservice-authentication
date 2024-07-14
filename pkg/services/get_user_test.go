package services_test

import (
	"context"
	"github.com/in-rich/uservice-authentication/config"
	"github.com/in-rich/uservice-authentication/pkg/dao"
	daomocks "github.com/in-rich/uservice-authentication/pkg/dao/mocks"
	"github.com/in-rich/uservice-authentication/pkg/entities"
	"github.com/in-rich/uservice-authentication/pkg/models"
	"github.com/in-rich/uservice-authentication/pkg/services"
	"github.com/stretchr/testify/require"
	"testing"
)

var getUserInfoFixtures = []*FixtureUser{
	{
		Email:         "user@gmail.com",
		EmailVerified: true,
		DisplayName:   "user one",
		UID:           "user-one-uid",
		Password:      "password",
		PhotoURL:      "https://image.png",
	},
}

func TestGetUserService(t *testing.T) {
	testData := []struct {
		name string

		uid string

		shouldCallGetUser bool
		getUserResponse   *entities.User
		getUserErr        error

		expect    *models.User
		expectErr error
	}{
		{
			name:              "GetUserInfo",
			uid:               "user-one-uid",
			shouldCallGetUser: true,
			getUserResponse: &entities.User{
				PublicIdentifier: "public-identifier-1",
				FirebaseUID:      "user-one-uid",
			},
			expect: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirebaseUID:      "user-one-uid",
				Email:            "user@gmail.com",
			},
		},
		{
			name:              "NoLocalData",
			uid:               "user-one-uid",
			shouldCallGetUser: true,
			getUserErr:        dao.ErrUserNotFound,
			expect: &models.User{
				PublicIdentifier: "",
				FirebaseUID:      "user-one-uid",
				Email:            "user@gmail.com",
			},
		},
		{
			name:      "UserNotFound",
			uid:       "user-two-uid",
			expectErr: services.ErrUserNotFound,
		},
		{
			name:              "GetUserError",
			uid:               "user-one-uid",
			shouldCallGetUser: true,
			getUserErr:        FooErr,
			expectErr:         FooErr,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			require.NoError(t, CreateUsersFixtures(getUserInfoFixtures))
			defer CleanUsersFixtures(getUserInfoFixtures)

			getUserRepository := daomocks.NewMockGetUserRepository(t)
			if data.shouldCallGetUser {
				getUserRepository.On("GetUser", context.TODO(), data.uid).Return(data.getUserResponse, data.getUserErr)
			}

			service := services.NewGetUserService(config.AuthClient, getUserRepository)

			user, err := service.Exec(context.TODO(), data.uid)

			require.ErrorIs(t, err, data.expectErr)
			require.Equal(t, data.expect, user)

			getUserRepository.AssertExpectations(t)
		})
	}
}
