package postgresTestSetup

import (
	"testing"
)

func Test_should_load_testcontainers(t *testing.T) {
	db := StartTestcontainersPostgres(t, "")

	defer teardownTestcontainersPostgres(t)

	if db == nil {
		t.Fatalf("db is nil")
	}
	db.Close()
}
