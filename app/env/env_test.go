package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var envValues = []Name{
	DBName,
	DBUser,
	DBPassword,
	DBHost,
	DBPort,
	MigrationDir,
	AdminUser,
	AdminPassword,
	OAuth2Domain,
	OAuth2ClientId,
	OAuth2ClientSecret,
	OAuth2Database,
}

func Test_should_load_env_file(t *testing.T) {
	err := Load(".env")
	assert.Nil(t, err)
	for _, envValue := range envValues {

		value := envValue.GetValue()
		assert.NotEmpty(t, value, "env value %s is empty", envValue)
	}
}
