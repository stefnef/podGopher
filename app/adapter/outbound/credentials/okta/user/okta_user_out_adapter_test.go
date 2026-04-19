package user

import (
	"context"
	"podGopher/core/domain/model"
	"podGopher/env"
	"testing"
	"time"

	"github.com/auth0/go-auth0/v2/management"
	"github.com/auth0/go-auth0/v2/management/client"
	"github.com/auth0/go-auth0/v2/management/option"
	"github.com/stretchr/testify/assert"
)

var (
	authMgmt    *client.Management
	exampleUser = &model.User{
		Id:        "123-456",
		Username:  "unit-test-user",
		IsAdmin:   false,
		ShowRoles: nil,
	}
	email    = "unit-test-user@example.com"
	password = "Ü+some8Very7Secret6Phrase/#1*?"
)

var oktaAuth *OktaUserOutAdapter

func setup(t *testing.T) context.Context {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var err error

	err = env.Load("./../../../../../env/.env")
	assert.Nil(t, err)

	oktaAuth = NewOktaAuthCredentialService(env.OAuth2Domain.GetValue(), env.OAuth2ClientId.GetValue(), env.OAuth2ClientSecret.GetValue())

	authMgmt, err = client.New(
		env.OAuth2Domain.GetValue(),
		option.WithClientCredentials(
			context.Background(),
			env.OAuth2ClientId.GetValue(),
			env.OAuth2ClientSecret.GetValue(),
		),
	)
	assert.Nil(t, err)
	return context.Background()
}

func searchUser(t *testing.T, ctx context.Context) []*management.UserResponseSchema {
	searchRequest := &management.ListUsersRequestParameters{
		Q:       management.String("username:" + exampleUser.Username),
		Sort:    management.String("created_at:1"),
		PerPage: management.Int(2),
	}

	usersPage, err := authMgmt.Users.List(ctx, searchRequest)
	assert.Nil(t, err)

	return usersPage.Results
}

func deleteUser(t *testing.T, ctx context.Context, userId string) {
	err := authMgmt.Users.Delete(ctx, userId)
	if err != nil {
		assert.Nil(t, err)
	}
}

func teardown(t *testing.T, ctx context.Context) {
	users := searchUser(t, ctx)
	for _, createdUser := range users {
		deleteUser(t, ctx, *createdUser.UserID)
	}
}

func Test_should_fail_with_invalid_credentials_on_New(t *testing.T) {
	service := NewOktaAuthCredentialService("foo", "bar", "")
	err := service.CreateUserCredentials(exampleUser, email, password)
	assert.NotNil(t, err)
}

func TestOktaUserOutAdapter_CreateUserCredentials(t *testing.T) {
	ctx := setup(t)
	defer teardown(t, ctx)

	t.Run("should not contain test users at beginning", func(t *testing.T) {
		users := searchUser(t, ctx)
		assert.Len(t, users, 0)
	})

	t.Run("should create user with valid credentials", func(t *testing.T) {

		err := oktaAuth.CreateUserCredentials(exampleUser, email, password)
		assert.Nil(t, err)

		time.Sleep(2000 * time.Millisecond)

		users := searchUser(t, ctx)
		assert.Len(t, users, 1)
	})

}
