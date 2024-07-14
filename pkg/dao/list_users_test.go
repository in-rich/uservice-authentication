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

var listUserFixtures = []*entities.User{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		PublicIdentifier: "public-identifier-1",
		FirebaseUID:      "firebase-uid-1",
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
		PublicIdentifier: "public-identifier-2",
		FirebaseUID:      "firebase-uid-2",
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
		PublicIdentifier: "public-identifier-3",
		FirebaseUID:      "firebase-uid-3",
	},
}

func TestListUsers(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name         string
		firebaseUIDs []string
		expect       []*entities.User
	}{
		{
			name:         "ListUsers",
			firebaseUIDs: []string{"firebase-uid-1", "firebase-uid-3", "firebase-uid-4"},
			expect: []*entities.User{
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					PublicIdentifier: "public-identifier-1",
					FirebaseUID:      "firebase-uid-1",
				},
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
					PublicIdentifier: "public-identifier-3",
					FirebaseUID:      "firebase-uid-3",
				},
			},
		},
		{
			name:         "ListUsersEmpty",
			firebaseUIDs: []string{"firebase-uid-4"},
			expect:       []*entities.User{},
		},
	}

	stx := BeginTX(db, listUserFixtures)
	defer RollbackTX(stx)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewListUsersRepository(tx)
			users, err := repo.ListUsers(context.TODO(), data.firebaseUIDs)

			require.NoError(t, err)
			require.Equal(t, data.expect, users)
		})
	}
}
