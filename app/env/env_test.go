package env

import (
	"os"
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
	CredentialService,
	OAuth2Domain,
	OAuth2ClientId,
	OAuth2ClientSecret,
	OAuth2Database,
}

var envFiles = []string{".env", ".example-env"}

func Test_should_load_env_files(t *testing.T) {
	for _, envFile := range envFiles {
		t.Run(envFile, func(t *testing.T) {
			err := Load(".env")
			assert.Nil(t, err)
			testEnvFile(t)
		})
	}
}

func testEnvFile(t *testing.T) {
	for _, envValue := range envValues {
		value := envValue.GetValue()
		assert.NotEmpty(t, value, "env value %s is empty", envValue)
		resetEnvVariable(t, envValue)
	}
}

func resetEnvVariable(t *testing.T, envValue Name) {
	err := os.Unsetenv(string(envValue))
	assert.Nil(t, err)
}
