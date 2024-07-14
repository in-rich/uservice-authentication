package services_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/in-rich/uservice-authentication/config"
	"github.com/in-rich/uservice-authentication/pkg/dao"
	daomocks "github.com/in-rich/uservice-authentication/pkg/dao/mocks"
	"github.com/in-rich/uservice-authentication/pkg/entities"
	"github.com/in-rich/uservice-authentication/pkg/models"
	"github.com/in-rich/uservice-authentication/pkg/services"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var authenticateFixtures = []*FixtureUser{
	{
		Email:         "user@gmail.com",
		EmailVerified: true,
		DisplayName:   "user one",
		UID:           "user-one-uid",
		Password:      "password",
		PhotoURL:      "https://image.png",
	},
	{
		Email:         "user2@gmail.com",
		EmailVerified: false,
		DisplayName:   "user rwo",
		UID:           "user-two-uid",
		Password:      "password",
		PhotoURL:      "https://image.png",
	},
}

func getIDToken(t *testing.T, uid string) string {
	const verifyTokenURL = "http://127.0.0.1:1151/identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken"

	// https://stackoverflow.com/questions/48268478/in-firebase-how-to-generate-an-idtoken-on-the-server-for-testing-purposes
	customToken, err := config.AuthClient.CustomToken(context.TODO(), uid)
	require.NoError(t, err)

	jsonData, err := json.Marshal(map[string]interface{}{
		"token":             customToken,
		"returnSecureToken": true,
	})
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s?key=%s", verifyTokenURL, config.Firebase.APIKey),
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		strBody := new(bytes.Buffer)
		_, _ = strBody.ReadFrom(resp.Body)
		t.Fatalf("sign in failed with status %s: %s", resp.Status, strBody.String())
	}

	var responseData struct {
		IDToken string `json:"idToken"`
	}

	err = json.NewDecoder(resp.Body).Decode(&responseData)
	require.NoError(t, err)

	return responseData.IDToken
}

func TestAuthenticate(t *testing.T) {
	require.NoError(t, CreateUsersFixtures(authenticateFixtures))
	defer CleanUsersFixtures(authenticateFixtures)

	validIDToken := getIDToken(t, "user-one-uid")
	emailNotValidatedIDToken := getIDToken(t, "user-two-uid")

	testData := []struct {
		name string

		token string

		shouldCallGetUser bool
		getUserResponse   *entities.User
		getUserErr        error

		expect    *models.User
		expectErr error
	}{
		{
			name:              "ValidToken",
			token:             validIDToken,
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
			name:              "NoExtraData",
			token:             validIDToken,
			shouldCallGetUser: true,
			getUserErr:        dao.ErrUserNotFound,
			expect: &models.User{
				PublicIdentifier: "",
				FirebaseUID:      "user-one-uid",
				Email:            "user@gmail.com",
			},
		},
		{
			name:              "GetUserError",
			token:             validIDToken,
			shouldCallGetUser: true,
			getUserErr:        FooErr,
			expectErr:         FooErr,
		},
		{
			name:      "EmptyToken",
			token:     "",
			expectErr: services.ErrUnauthenticated,
		},
		{
			name:      "InvalidToken",
			token:     "invalid-token",
			expectErr: services.ErrVerifyToken,
		},
		{
			name:      "EmailNotVerified",
			token:     emailNotValidatedIDToken,
			expectErr: services.ErrEmailNotVerified,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			getUserRepository := daomocks.NewMockGetUserRepository(t)

			if tt.shouldCallGetUser {
				getUserRepository.On("GetUser", context.TODO(), "user-one-uid").Return(tt.getUserResponse, tt.getUserErr)
			}

			service := services.NewAuthenticateService(config.AuthClient, getUserRepository)

			user, err := service.Exec(context.TODO(), tt.token)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, user)

			getUserRepository.AssertExpectations(t)
		})
	}
}
