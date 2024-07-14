package dao

import (
	"context"
	"github.com/in-rich/uservice-authentication/pkg/entities"
	"github.com/uptrace/bun"
)

type UpdateUserData struct {
	PublicIdentifier string
}

type UpdateUserRepository interface {
	UpdateUser(ctx context.Context, firebaseUID string, data *UpdateUserData) (*entities.User, error)
}

type updateUserRepositoryImpl struct {
	db bun.IDB
}

func (r *updateUserRepositoryImpl) UpdateUser(ctx context.Context, firebaseUID string, data *UpdateUserData) (*entities.User, error) {
	user := &entities.User{
		PublicIdentifier: data.PublicIdentifier,
		FirebaseUID:      firebaseUID,
	}

	res, err := r.db.NewUpdate().
		Model(user).
		Column("public_identifier").
		Where("firebase_uid = ?", firebaseUID).
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func NewUpdateUserRepository(db bun.IDB) UpdateUserRepository {
	return &updateUserRepositoryImpl{
		db: db,
	}
}
