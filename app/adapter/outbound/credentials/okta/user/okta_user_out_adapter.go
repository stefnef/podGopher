package user

import (
	"context"
	"podGopher/core/domain/model"

	"github.com/auth0/go-auth0/v2/management"
	"github.com/auth0/go-auth0/v2/management/client"
	"github.com/auth0/go-auth0/v2/management/option"
)

type OktaUserOutAdapter struct {
	mgmt *client.Management
}

func NewOktaAuthCredentialService(domain string, clientId string, clientSecret string) *OktaUserOutAdapter {
	mgmt, _ := client.New(
		domain,
		option.WithClientCredentials(
			context.Background(),
			clientId,
			clientSecret,
		),
	)

	return &OktaUserOutAdapter{mgmt: mgmt}
}

func (okta *OktaUserOutAdapter) CreateUserCredentials(newUser *model.User, email string, password string) (_ error) {
	ctx := context.Background()
	createUserRequest := &management.CreateUserRequestContent{
		Email:      management.String(email),
		Username:   management.String(newUser.Username),
		Connection: "Username-Password-Authentication",
		Password:   management.String(password),
		UserMetadata: &management.UserMetadata{
			"preference": "email",
		},
	}

	_, err := okta.mgmt.Users.Create(ctx, createUserRequest)
	return err
}
