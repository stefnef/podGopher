package episode

import (
	"database/sql"
	"podGopher/core/domain/model"
)

type PostgresEpisodeOutAdapter struct {
	db *sql.DB
}

func (adapter *PostgresEpisodeOutAdapter) SaveEpisode(episode *model.Episode) (err error) {
	if err = adapter.createEpisodeEntry(episode); err != nil {
		return err
	}
	if err = adapter.createShowEpisodeMappingEntry(episode); err != nil {
		return err //TODO clean up on error -> transactional
	}
	return nil
}

func (adapter *PostgresEpisodeOutAdapter) createShowEpisodeMappingEntry(episode *model.Episode) (err error) {
	var stmt *sql.Stmt

	if stmt, err = adapter.db.Prepare("INSERT INTO show_episodes (show_id, episode_id) VALUES ($1, $2);"); err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	if _, err = stmt.Exec(episode.ShowId, episode.Id); err != nil {
		return err
	}

	return nil
}

func (adapter *PostgresEpisodeOutAdapter) createEpisodeEntry(episode *model.Episode) (err error) {
	var stmt *sql.Stmt

	if stmt, err = adapter.db.Prepare("INSERT INTO episode (id, show_id, title) VALUES ($1, $2, $3);"); err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	if _, err = stmt.Exec(episode.Id, episode.ShowId, episode.Title); err != nil {
		return err
	}

	return nil
}

func (adapter *PostgresEpisodeOutAdapter) ExistsByTitle(title string) bool {
	query := "SELECT EXISTS(SELECT 1 FROM episode where title = $1)"
	row := adapter.db.QueryRow(query, title)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (adapter *PostgresEpisodeOutAdapter) GetEpisodeOrNil(id string) (episode *model.Episode, err error) {
	query := "SELECT * FROM episode where id = $1"
	row := adapter.db.QueryRow(query, id)

	episode = &model.Episode{}
	if err = row.Scan(&episode.Id, &episode.ShowId, &episode.Title); err != nil {
		return nil, nil
	}
	return episode, nil
}

func NewPostgresEpisodeRepository(db *sql.DB) *PostgresEpisodeOutAdapter {
	return &PostgresEpisodeOutAdapter{db: db}
}
