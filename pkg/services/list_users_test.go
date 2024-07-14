package services_test

import (
	"context"
	"github.com/in-rich/uservice-authentication/config"
	daomocks "github.com/in-rich/uservice-authentication/pkg/dao/mocks"
	"github.com/in-rich/uservice-authentication/pkg/entities"
	"github.com/in-rich/uservice-authentication/pkg/models"
	"github.com/in-rich/uservice-authentication/pkg/services"
	"github.com/stretchr/testify/require"
	"testing"
)

var listUsersInfoFixtures = []*FixtureUser{
	{
		Email:         "user1@gmail.com",
		EmailVerified: true,
		DisplayName:   "user one",
		UID:           "user-one-uid",
		Password:      "password",
		PhotoURL:      "https://image.png",
	},
	{
		Email:         "user2@gmail.com",
		EmailVerified: true,
		DisplayName:   "user two",
		UID:           "user-two-uid",
		Password:      "password",
		PhotoURL:      "https://image.png",
	},
	{
		Email:         "user3@gmail.com",
		EmailVerified: true,
		DisplayName:   "user three",
		UID:           "user-three-uid",
		Password:      "password",
		PhotoURL:      "https://image.png",
	},
}

func TestListUsersService(t *testing.T) {
	testData := []struct {
		name string

		uids []string

		listUsersResult []*entities.User
		listUsersErr    error

		expect    []*models.User
		expectErr error
	}{
		{
			name: "ListUsers",
			uids: []string{"user-one-uid", "user-four-uid", "user-three-uid"},
			listUsersResult: []*entities.User{
				{
					PublicIdentifier: "public-identifier-1",
					FirebaseUID:      "user-one-uid",
				},
			},
			expect: []*models.User{
				{
					PublicIdentifier: "public-identifier-1",
					FirebaseUID:      "user-one-uid",
					Email:            "user1@gmail.com",
				},
				{
					PublicIdentifier: "",
					FirebaseUID:      "user-three-uid",
					Email:            "user3@gmail.com",
				},
			},
		},
		{
			name:         "ListUsersError",
			uids:         []string{"user-one-uid", "user-four-uid", "user-three-uid"},
			listUsersErr: FooErr,
			expectErr:    FooErr,
		},
		{
			name:            "NoResults",
			uids:            []string{"user-four-uid"},
			listUsersResult: []*entities.User{},
			expect:          []*models.User{},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			require.NoError(t, CreateUsersFixtures(listUsersInfoFixtures))
			defer CleanUsersFixtures(listUsersInfoFixtures)

			listUsersRepository := daomocks.NewMockListUsersRepository(t)
			listUsersRepository.On("ListUsers", context.TODO(), data.uids).
				Return(data.listUsersResult, data.listUsersErr)

			service := services.NewListUsersService(config.AuthClient, listUsersRepository)

			users, err := service.Exec(context.TODO(), data.uids)

			require.ErrorIs(t, err, data.expectErr)
			require.Equal(t, data.expect, users)

			listUsersRepository.AssertExpectations(t)
		})
	}
}
