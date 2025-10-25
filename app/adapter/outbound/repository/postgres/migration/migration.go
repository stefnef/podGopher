package migration

import (
	"fmt"
	"podGopher/env"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migration struct {
	migrate *migrate.Migrate
}

func NewMigration() (*Migration, error) {
	dsn := GetPostgresConnectionString()
	sourceUrl := fmt.Sprintf("file://%s", env.MigrationDir.GetValue())

	m, err := migrate.New(sourceUrl, dsn)
	if err != nil {
		return nil, err
	}
	return &Migration{m}, err
}

func GetPostgresConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.DBUser.GetValue(),
		env.DBPassword.GetValue(),
		env.DBHost.GetValue(),
		env.DBPort.GetValue(),
		env.DBName.GetValue(),
	)
}

func (m *Migration) Migrate() error {
	defer func(migrate *migrate.Migrate) {
		_, _ = migrate.Close()
	}(m.migrate)

	err := m.migrate.Up()
	return err
}
