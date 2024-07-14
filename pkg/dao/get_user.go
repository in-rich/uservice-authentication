package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/in-rich/uservice-authentication/pkg/entities"
	"github.com/uptrace/bun"
)

type GetUserRepository interface {
	GetUser(ctx context.Context, uid string) (*entities.User, error)
}

type getUserRepositoryImpl struct {
	db bun.IDB
}

func (r *getUserRepositoryImpl) GetUser(ctx context.Context, uid string) (*entities.User, error) {
	user := new(entities.User)

	err := r.db.NewSelect().Model(user).Where("firebase_uid = ?", uid).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func NewGetUserRepository(db bun.IDB) GetUserRepository {
	return &getUserRepositoryImpl{
		db: db,
	}
}
