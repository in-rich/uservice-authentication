package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-authentication/pkg/dao"
	"github.com/in-rich/uservice-authentication/pkg/models"
)

type UpdateUserService interface {
	Exec(ctx context.Context, token string, user *models.UpdateUser) (*models.User, error)
}

type updateUserServiceImpl struct {
	auth      AuthenticateService
	createDAO dao.CreateUserRepository
	updateDAO dao.UpdateUserRepository
}

func (s *updateUserServiceImpl) Exec(ctx context.Context, token string, data *models.UpdateUser) (*models.User, error) {
	firebaseUser, err := s.auth.Exec(ctx, token)
	if err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidUpdateUser, err)
	}

	user, err := s.createDAO.CreateUser(ctx, firebaseUser.FirebaseUID, &dao.CreateUserData{
		PublicIdentifier: data.PublicIdentifier,
	})

	// User was successfully created.
	if err == nil {
		return &models.User{
			PublicIdentifier: user.PublicIdentifier,
			FirebaseUID:      firebaseUser.FirebaseUID,
			Email:            firebaseUser.Email,
		}, nil
	}

	if !errors.Is(err, dao.ErrUserAlreadyExists) {
		return nil, err
	}

	// User already existed. Update it.
	user, err = s.updateDAO.UpdateUser(ctx, firebaseUser.FirebaseUID, &dao.UpdateUserData{
		PublicIdentifier: data.PublicIdentifier,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		PublicIdentifier: user.PublicIdentifier,
		FirebaseUID:      firebaseUser.FirebaseUID,
		Email:            firebaseUser.Email,
	}, nil
}

func NewUpdateUserService(
	auth AuthenticateService,
	createDAO dao.CreateUserRepository,
	updateDAO dao.UpdateUserRepository,
) UpdateUserService {
	return &updateUserServiceImpl{
		auth:      auth,
		createDAO: createDAO,
		updateDAO: updateDAO,
	}
}
