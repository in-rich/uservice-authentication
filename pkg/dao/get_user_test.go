package dao_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-authentication/pkg/dao"
	"github.com/in-rich/uservice-authentication/pkg/entities"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
)

var getUserFixtures = []*entities.User{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		PublicIdentifier: "public-identifier-1",
		FirebaseUID:      "firebase-uid-1",
	},
}

func TestGetUser(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name        string
		firebaseUID string
		expect      *entities.User
		expectErr   error
	}{
		{
			name:        "GetUser",
			firebaseUID: "firebase-uid-1",
			expect: &entities.User{
				ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				PublicIdentifier: "public-identifier-1",
				FirebaseUID:      "firebase-uid-1",
			},
		},
		{
			name:        "UserNotFound",
			firebaseUID: "firebase-uid-2",
			expectErr:   dao.ErrUserNotFound,
		},
	}

	stx := BeginTX(db, getUserFixtures)
	defer RollbackTX(stx)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewGetUserRepository(tx)
			user, err := repo.GetUser(context.TODO(), data.firebaseUID)

			require.ErrorIs(t, err, data.expectErr)
			require.Equal(t, data.expect, user)
		})
	}
}
