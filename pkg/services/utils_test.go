package services_test

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/in-rich/uservice-authentication/config"
)

var FooErr = errors.New("foo error")

type FixtureUser struct {
	Email         string
	EmailVerified bool
	DisplayName   string
	UID           string
	Password      string
	PhotoURL      string
}

func CreateUsersFixtures(fixtures []*FixtureUser) error {
	// Just in case something went wrong on latest run.
	CleanUsersFixtures(getUserInfoFixtures)

	for _, fixture := range fixtures {
		user := new(auth.UserToCreate).
			Email(fixture.Email).
			EmailVerified(fixture.EmailVerified).
			DisplayName(fixture.DisplayName).
			UID(fixture.UID).
			Password(fixture.Password).
			PhotoURL(fixture.PhotoURL)

		_, err := config.AuthClient.CreateUser(context.TODO(), user)
		if err != nil {
			return err
		}
	}

	return nil
}

func CleanUsersFixtures(fixtures []*FixtureUser) {
	for _, fixture := range fixtures {
		_ = config.AuthClient.DeleteUser(context.TODO(), fixture.UID)
	}
}
