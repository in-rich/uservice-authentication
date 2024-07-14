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

var updateUserFixtures = []*entities.User{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		PublicIdentifier: "public-identifier-1",
		FirebaseUID:      "firebase-uid-1",
	},
}

func TestUpdateUser(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name        string
		firebaseUID string
		data        *dao.UpdateUserData
		expect      *entities.User
		expectErr   error
	}{
		{
			name:        "UpdateUser",
			firebaseUID: "firebase-uid-1",
			data: &dao.UpdateUserData{
				PublicIdentifier: "public-identifier-2",
			},
			expect: &entities.User{
				ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				PublicIdentifier: "public-identifier-2",
				FirebaseUID:      "firebase-uid-1",
			},
		},
		{
			name:        "UserNotFound",
			firebaseUID: "firebase-uid-2",
			data: &dao.UpdateUserData{
				PublicIdentifier: "public-identifier-2",
			},
			expectErr: dao.ErrUserNotFound,
		},
	}

	stx := BeginTX(db, updateUserFixtures)
	defer RollbackTX(stx)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewUpdateUserRepository(tx)
			user, err := repo.UpdateUser(context.TODO(), data.firebaseUID, data.data)

			require.ErrorIs(t, err, data.expectErr)
			require.Equal(t, data.expect, user)
		})
	}
}
