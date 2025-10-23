package postgresTestSetup

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"podGopher/adapter/outbound/repository/postgres/migration"
	"podGopher/env"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	postgresClient "gocloud.dev/postgres"
)

var postgresContainer *postgres.PostgresContainer

func TeardownTestcontainersPostgres(t *testing.T) {
	if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}

func StartTestcontainersPostgres(t *testing.T, configDir string) *sql.DB {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ctx := context.Background()

	createContainer(t, configDir, ctx)
	dsn := createConnectionUrl(t, ctx)

	db := createAndVerifyConnection(t, ctx, dsn)

	executeMigration(t, configDir)

	return db
}

func executeMigration(t *testing.T, configDir string) {
	if err := os.Setenv(string(env.MigrationDir), configDir+"../migration/files"); err != nil {
		t.Fatal(err)
	}

	m, err := migration.NewMigration()
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Migrate(); err != nil {
		t.Fatal(err)
	}
}

func createAndVerifyConnection(t *testing.T, ctx context.Context, dsn string) *sql.DB {
	db, err := postgresClient.Open(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}
	return db
}

func createConnectionUrl(t *testing.T, ctx context.Context) string {
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv(string(env.DBHost), host); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv(string(env.DBPort), mappedPort.Port()); err != nil {
		t.Fatal(err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.DBUser.GetValue(),
		env.DBPassword.GetValue(),
		env.DBHost.GetValue(),
		env.DBPort.GetValue(),
		env.DBName.GetValue(),
	)
	return dsn
}

func createContainer(t *testing.T, configDir string, ctx context.Context) {
	var err error

	if err := env.Load(configDir + "../../../../../env/.testcontainers-env"); err != nil {
		t.Fatal(err)
	}

	postgresContainer, err = postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithOrderedInitScripts(
			filepath.Join(configDir, "postgres_init.sh"),
		),
		postgres.WithDatabase(env.DBName.GetValue()),
		postgres.WithUsername(env.DBUser.GetValue()),
		postgres.WithPassword(env.DBPassword.GetValue()),
		postgres.BasicWaitStrategies(),
	)
	if err != nil || postgresContainer == nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}
}
