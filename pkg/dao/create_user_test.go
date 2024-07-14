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

var createUserFixtures = []*entities.User{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		PublicIdentifier: "public-identifier-1",
		FirebaseUID:      "firebase-uid-1",
	},
}

func TestCreateUser(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name        string
		firebaseUID string
		data        *dao.CreateUserData
		expect      *entities.User
		expectErr   error
	}{
		{
			name:        "CreateUser",
			firebaseUID: "firebase-uid-2",
			data: &dao.CreateUserData{
				PublicIdentifier: "public-identifier-1",
			},
			expect: &entities.User{
				PublicIdentifier: "public-identifier-1",
				FirebaseUID:      "firebase-uid-2",
			},
		},
		{
			name:        "UserAlreadyExists",
			firebaseUID: "firebase-uid-1",
			data: &dao.CreateUserData{
				PublicIdentifier: "public-identifier-1",
			},
			expectErr: dao.ErrUserAlreadyExists,
		},
	}

	stx := BeginTX(db, createUserFixtures)
	defer RollbackTX(stx)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewCreateUserRepository(tx)
			user, err := repo.CreateUser(context.TODO(), data.firebaseUID, data.data)

			if user != nil {
				// Since ID is random, nullify it for comparison.
				user.ID = nil
			}

			require.ErrorIs(t, err, data.expectErr)
			require.Equal(t, data.expect, user)
		})
	}
}
