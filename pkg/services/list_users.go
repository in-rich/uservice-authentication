package services

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/in-rich/uservice-authentication/pkg/dao"
	"github.com/in-rich/uservice-authentication/pkg/entities"
	"github.com/in-rich/uservice-authentication/pkg/models"
	"github.com/samber/lo"
)

type ListUsersService interface {
	Exec(ctx context.Context, uids []string) ([]*models.User, error)
}

type listUsersServiceImpl struct {
	client *auth.Client
	dao    dao.ListUsersRepository
}

func (s *listUsersServiceImpl) Exec(ctx context.Context, uids []string) ([]*models.User, error) {
	identifiers := lo.Map(uids, func(item string, index int) auth.UserIdentifier {
		return auth.UIDIdentifier{UID: item}
	})

	users, err := s.client.GetUsers(ctx, identifiers)
	if err != nil {
		return nil, err
	}

	extras, err := s.dao.ListUsers(ctx, uids)
	if err != nil {
		return nil, err
	}

	return lo.Map(users.Users, func(item *auth.UserRecord, index int) *models.User {
		extra, ok := lo.Find(extras, func(extraItem *entities.User) bool {
			return extraItem.FirebaseUID == item.UID
		})

		result := &models.User{
			FirebaseUID: item.UID,
			Email:       item.Email,
		}

		if ok {
			result.PublicIdentifier = extra.PublicIdentifier
		}

		return result
	}), nil
}

func NewListUsersService(client *auth.Client, dao dao.ListUsersRepository) ListUsersService {
	return &listUsersServiceImpl{
		client: client,
		dao:    dao,
	}
}
