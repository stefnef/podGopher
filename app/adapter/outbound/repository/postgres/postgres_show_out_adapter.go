package repository

import (
	"database/sql"
	"podGopher/core/port/outbound"

	"github.com/google/uuid"
)

type PostgresShowOutAdapter struct {
	db *sql.DB
}

func (adapter *PostgresShowOutAdapter) SaveShow(title string) (err error) {
	var stmt *sql.Stmt
	id := uuid.NewString()

	if stmt, err = adapter.db.Prepare("INSERT INTO shows (id, title) VALUES ($1, $2);"); err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	if _, err = stmt.Exec(id, title); err != nil {
		return err
	}

	return nil
}

func (adapter *PostgresShowOutAdapter) ExistsByTitle(title string) bool {
	query := "SELECT EXISTS(SELECT 1 FROM shows where title = $1)"
	row := adapter.db.QueryRow(query, title)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func NewPostgresShowRepository(db *sql.DB) outbound.SaveShowPort {
	return &PostgresShowOutAdapter{db: db}
}
