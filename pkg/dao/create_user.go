package dao

import (
	"context"
	"errors"
	"github.com/in-rich/uservice-authentication/pkg/entities"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type CreateUserData struct {
	PublicIdentifier string
}

type CreateUserRepository interface {
	CreateUser(ctx context.Context, firebaseUID string, data *CreateUserData) (*entities.User, error)
}

type createUserRepositoryImpl struct {
	db bun.IDB
}

func (r *createUserRepositoryImpl) CreateUser(ctx context.Context, firebaseUID string, data *CreateUserData) (*entities.User, error) {
	user := &entities.User{
		PublicIdentifier: data.PublicIdentifier,
		FirebaseUID:      firebaseUID,
	}

	if _, err := r.db.NewInsert().Model(user).Returning("*").Exec(ctx); err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() {
			return nil, ErrUserAlreadyExists
		}

		return nil, err
	}

	return user, nil
}

func NewCreateUserRepository(db bun.IDB) CreateUserRepository {
	return &createUserRepositoryImpl{
		db: db,
	}
}
