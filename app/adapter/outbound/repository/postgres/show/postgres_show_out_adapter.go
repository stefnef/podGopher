package show

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
	query := "SELECT s.*, se.episode_id FROM show s LEFT JOIN show_episodes se ON se.show_id = s.id WHERE s.id = $1;"
	rows, _ := adapter.db.Query(query, id)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		show, err = parseNextShow(rows, show)
		if err != nil {
			return nil, err
		}
	}

	return show, nil
}

func parseNextShow(rows *sql.Rows, show *model.Show) (*model.Show, error) {
	var (
		showId string
		title  string
		slug   string
		eId    sql.NullString
	)

	if err := rows.Scan(&showId, &title, &slug, &eId); err != nil {
		return nil, err
	}

	if show == nil {
		show = &model.Show{
			Id:    showId,
			Title: title,
			Slug:  slug,
		}
	}

	if eId.Valid {
		show.Episodes = append(show.Episodes, eId.String)
	}
	return show, nil
}

func NewPostgresShowRepository(db *sql.DB) *PostgresShowOutAdapter {
	return &PostgresShowOutAdapter{db: db}
}
