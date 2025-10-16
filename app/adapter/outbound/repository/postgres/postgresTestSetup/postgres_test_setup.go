package postgresTestSetup

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	postgresClient "gocloud.dev/postgres"
)

var testcontainersDbConfig = struct {
	dbName     string
	dbUser     string
	dbPassword string
}{
	dbName:     "podcasts",
	dbUser:     "user",
	dbPassword: "password",
}

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

	return db
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
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		testcontainersDbConfig.dbUser,
		testcontainersDbConfig.dbPassword,
		host,
		mappedPort.Port(),
		testcontainersDbConfig.dbName,
	)
	return dsn
}

func createContainer(t *testing.T, configDir string, ctx context.Context) {
	var err error
	postgresContainer, err = postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithOrderedInitScripts(
			filepath.Join(configDir, "postgres_init.sh"),
			filepath.Join(configDir, "../setup/001_init_shows.sql"),
		),
		postgres.WithDatabase(testcontainersDbConfig.dbName),
		postgres.WithUsername(testcontainersDbConfig.dbUser),
		postgres.WithPassword(testcontainersDbConfig.dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil || postgresContainer == nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}
}
