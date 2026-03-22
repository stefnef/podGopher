package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_should_create_Auth(t *testing.T) {
	auth := NewAdminAuth("fake-user", "fake-password")

	assert.NotNil(t, auth)
}

func Test_should_verify_credentials_for_AdminAuth(t *testing.T) {
	auth := NewAdminAuth("fake-user", "fake-password")

	assert.True(t, auth.IsValid("fake-user", "fake-password"))
	assert.False(t, auth.IsValid("fake-user", "wrong-password"))
	assert.False(t, auth.IsValid("wrong-user", "fake-password"))

}

func Test_should_not_accept_empty_credentials_for_AdminAuth(t *testing.T) {
	assert.PanicsWithValue(t, "admin username or password cannot be empty", func() {
		NewAdminAuth("", "fake-password")
	})

	assert.Panics(t, func() {
		NewAdminAuth("fake-user", "")
	})

}
