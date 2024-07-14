package services

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/in-rich/uservice-authentication/pkg/dao"
	"github.com/in-rich/uservice-authentication/pkg/entities"
	"github.com/in-rich/uservice-authentication/pkg/models"
)

type GetUserService interface {
	Exec(ctx context.Context, uid string) (*models.User, error)
}

type getUserServiceImpl struct {
	client *auth.Client
	dao    dao.GetUserRepository
}

func (s *getUserServiceImpl) Exec(ctx context.Context, uid string) (*models.User, error) {
	user, err := s.client.GetUser(ctx, uid)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	extra, err := s.dao.GetUser(ctx, uid)
	if err != nil {
		if !errors.Is(err, dao.ErrUserNotFound) {
			return nil, err
		}

		// If no extra information found, just return the firebase user with their default value.
		extra = new(entities.User)
	}

	return &models.User{
		PublicIdentifier: extra.PublicIdentifier,
		FirebaseUID:      user.UID,
		Email:            user.Email,
	}, nil
}

func NewGetUserService(client *auth.Client, dao dao.GetUserRepository) GetUserService {
	return &getUserServiceImpl{
		client: client,
		dao:    dao,
	}
}
