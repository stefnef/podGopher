package repository

import (
	"database/sql"
	"podGopher/core/domain/model"
)

type PostgresShowOutAdapter struct {
	db *sql.DB
}

func (adapter *PostgresShowOutAdapter) SaveShow(show *model.Show) (err error) {
	var stmt *sql.Stmt

	if stmt, err = adapter.db.Prepare("INSERT INTO show (id, title, slug) VALUES ($1, $2, $3);"); err != nil {
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
	query := "SELECT EXISTS(SELECT 1 FROM show where title = $1 or slug = $2)"
	row := adapter.db.QueryRow(query, title, slug)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (adapter *PostgresShowOutAdapter) GetShowOrNil(id string) (show *model.Show, err error) {
	query := "SELECT * FROM show where id = $1"
	row := adapter.db.QueryRow(query, id)

	show = &model.Show{}
	if err = row.Scan(&show.Id, &show.Title, &show.Slug); err != nil {
		return nil, nil
	}
	return show, nil
}

func NewPostgresShowRepository(db *sql.DB) *PostgresShowOutAdapter {
	return &PostgresShowOutAdapter{db: db}
}
