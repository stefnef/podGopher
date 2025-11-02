package repository

import (
	"database/sql"
	"podGopher/core/domain/model"
	"podGopher/core/port/outbound"
)

type PostgresShowOutAdapter struct {
	db *sql.DB
}

func (adapter *PostgresShowOutAdapter) SaveShow(show *model.Show) (err error) {
	var stmt *sql.Stmt

	if stmt, err = adapter.db.Prepare("INSERT INTO shows (id, title, slug) VALUES ($1, $2, $3);"); err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	if _, err = stmt.Exec(show.Id, show.Title, show.Slug); err != nil {
		return err
	}

	return nil
}

func (adapter *PostgresShowOutAdapter) ExistsByTitleOrSlug(title string, slug string) bool {
	query := "SELECT EXISTS(SELECT 1 FROM shows where title = $1 or slug = $2)"
	row := adapter.db.QueryRow(query, title, slug)

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
