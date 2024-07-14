package services

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/in-rich/uservice-authentication/pkg/dao"
	"github.com/in-rich/uservice-authentication/pkg/entities"
	"github.com/in-rich/uservice-authentication/pkg/models"
)

type AuthenticateService interface {
	Exec(ctx context.Context, token string) (*models.User, error)
}

type authenticateServiceImpl struct {
	client            *auth.Client
	getUserRepository dao.GetUserRepository
}

func (s *authenticateServiceImpl) Exec(ctx context.Context, token string) (*models.User, error) {
	if token == "" {
		return nil, ErrUnauthenticated
	}

	authToken, err := s.client.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, errors.Join(ErrVerifyToken, err)
	}

	user, err := s.client.GetUser(ctx, authToken.UID)
	if err != nil {
		return nil, err
	}

	// Force users to verify their email to use the service.
	if user.EmailVerified == false {
		return nil, ErrEmailNotVerified
	}

	extra, err := s.getUserRepository.GetUser(ctx, user.UID)
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

func NewAuthenticateService(client *auth.Client, getUserRepository dao.GetUserRepository) AuthenticateService {
	return &authenticateServiceImpl{
		client:            client,
		getUserRepository: getUserRepository,
	}
}
