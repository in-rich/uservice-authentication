package services_test

import (
	"context"
	"github.com/in-rich/uservice-authentication/pkg/dao"
	daomocks "github.com/in-rich/uservice-authentication/pkg/dao/mocks"
	"github.com/in-rich/uservice-authentication/pkg/entities"
	"github.com/in-rich/uservice-authentication/pkg/models"
	"github.com/in-rich/uservice-authentication/pkg/services"
	servicesmocks "github.com/in-rich/uservice-authentication/pkg/services/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateUser(t *testing.T) {
	testData := []struct {
		name string

		token string
		data  *models.UpdateUser

		authResponse *models.User
		authErr      error

		shouldCallCreateUser bool
		createUserResponse   *entities.User
		createUserErr        error

		shouldCallUpdateUser bool
		updateUserResponse   *entities.User
		updateUserErr        error

		expect    *models.User
		expectErr error
	}{
		{
			name:  "CreateUser",
			token: "foo-token",
			data: &models.UpdateUser{
				PublicIdentifier: "public-identifier-2",
			},
			authResponse: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirebaseUID:      "user-one-uid",
				Email:            "user@gmail.com",
			},
			shouldCallCreateUser: true,
			createUserResponse: &entities.User{
				FirebaseUID:      "user-one-uid",
				PublicIdentifier: "public-identifier-2",
			},
			expect: &models.User{
				FirebaseUID:      "user-one-uid",
				PublicIdentifier: "public-identifier-2",
				Email:            "user@gmail.com",
			},
		},
		{
			name:  "UpdateUser",
			token: "foo-token",
			data: &models.UpdateUser{
				PublicIdentifier: "public-identifier-2",
			},
			authResponse: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirebaseUID:      "user-one-uid",
				Email:            "user@gmail.com",
			},
			shouldCallCreateUser: true,
			createUserErr:        dao.ErrUserAlreadyExists,
			shouldCallUpdateUser: true,
			updateUserResponse: &entities.User{
				FirebaseUID:      "user-one-uid",
				PublicIdentifier: "public-identifier-2",
			},
			expect: &models.User{
				FirebaseUID:      "user-one-uid",
				PublicIdentifier: "public-identifier-2",
				Email:            "user@gmail.com",
			},
		},
		{
			name:  "AuthError",
			token: "foo-token",
			data: &models.UpdateUser{
				PublicIdentifier: "public-identifier-2",
			},
			authErr:   FooErr,
			expectErr: FooErr,
		},
		{
			name:  "CreateUserError",
			token: "foo-token",
			data: &models.UpdateUser{
				PublicIdentifier: "public-identifier-2",
			},
			authResponse: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirebaseUID:      "user-one-uid",
				Email:            "user@gmail.com",
			},
			shouldCallCreateUser: true,
			createUserErr:        FooErr,
			expectErr:            FooErr,
		},
		{
			name:  "UpdateUserError",
			token: "foo-token",
			data: &models.UpdateUser{
				PublicIdentifier: "public-identifier-2",
			},
			authResponse: &models.User{
				PublicIdentifier: "public-identifier-1",
				FirebaseUID:      "user-one-uid",
				Email:            "user@gmail.com",
			},
			shouldCallCreateUser: true,
			createUserErr:        dao.ErrUserAlreadyExists,
			shouldCallUpdateUser: true,
			updateUserErr:        FooErr,
			expectErr:            FooErr,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			authService := servicesmocks.NewMockAuthenticateService(t)
			createUserRepository := daomocks.NewMockCreateUserRepository(t)
			updateUserRepository := daomocks.NewMockUpdateUserRepository(t)

			authService.On("Exec", context.TODO(), data.token).Return(data.authResponse, data.authErr)

			if data.shouldCallCreateUser {
				createUserRepository.On("CreateUser", context.TODO(), data.authResponse.FirebaseUID, &dao.CreateUserData{
					PublicIdentifier: data.data.PublicIdentifier,
				}).Return(data.createUserResponse, data.createUserErr)
			}

			if data.shouldCallUpdateUser {
				updateUserRepository.On("UpdateUser", context.TODO(), data.authResponse.FirebaseUID, &dao.UpdateUserData{
					PublicIdentifier: data.data.PublicIdentifier,
				}).Return(data.updateUserResponse, data.updateUserErr)
			}

			service := services.NewUpdateUserService(authService, createUserRepository, updateUserRepository)

			user, err := service.Exec(context.TODO(), data.token, data.data)

			require.ErrorIs(t, err, data.expectErr)
			require.Equal(t, data.expect, user)

			authService.AssertExpectations(t)
			createUserRepository.AssertExpectations(t)
			updateUserRepository.AssertExpectations(t)
		})
	}
}
