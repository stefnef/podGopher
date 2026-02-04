package distribution

import (
	"database/sql"
	"podGopher/core/domain/model"
)

type PostgresDistributionOutAdapter struct {
	db *sql.DB
}

func (adapter *PostgresDistributionOutAdapter) SaveDistribution(distribution *model.Distribution) (err error) {
	transaction, _ := adapter.db.Begin()
	defer func(transaction *sql.Tx) {
		_ = transaction.Rollback()
	}(transaction)

	if err = adapter.createDistributionEntry(distribution, transaction); err != nil {
		return err
	}
	if err = adapter.createShowDistributionMappingEntry(distribution, transaction); err != nil {
		return err
	}
	_ = transaction.Commit()
	return nil
}

func (adapter *PostgresDistributionOutAdapter) createShowDistributionMappingEntry(distribution *model.Distribution, transaction *sql.Tx) (err error) {
	var stmt *sql.Stmt

	if stmt, err = transaction.Prepare("INSERT INTO show_distributions (show_id, distribution_id) VALUES ($1, $2);"); err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	if _, err = stmt.Exec(distribution.ShowId, distribution.Id); err != nil {
		return err
	}

	return nil
}

func (adapter *PostgresDistributionOutAdapter) updateShowDistributionMappingEntry(distribution *model.Distribution, transaction *sql.Tx) (err error) {
	var stmt *sql.Stmt

	if stmt, err = transaction.Prepare("UPDATE show_distributions SET show_id = $1 WHERE distribution_id = $2;"); err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	if _, err = stmt.Exec(distribution.ShowId, distribution.Id); err != nil {
		return err
	}

	return nil
}

func (adapter *PostgresDistributionOutAdapter) deleteEpisodeDistributionMappingEntries(distribution *model.Distribution, transaction *sql.Tx) (err error) {
	var stmt *sql.Stmt

	if stmt, err = transaction.Prepare("DELETE FROM episodes_distributions WHERE distribution_id = $1;"); err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	if _, err = stmt.Exec(distribution.Id); err != nil {
		return err
	}

	return nil
}

func (adapter *PostgresDistributionOutAdapter) insertEpisodeDistributionMappingEntries(distribution *model.Distribution, transaction *sql.Tx) (err error) {
	var stmt *sql.Stmt

	if stmt, err = transaction.Prepare("INSERT INTO episodes_distributions (episode_id, distribution_id) VALUES ($1, $2);"); err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	for _, episodeId := range distribution.Episodes {
		if _, err = stmt.Exec(episodeId, distribution.Id); err != nil {
			return err
		}
	}

	return nil
}

func (adapter *PostgresDistributionOutAdapter) createDistributionEntry(distribution *model.Distribution, transaction *sql.Tx) (err error) {
	var stmt *sql.Stmt

	if stmt, err = transaction.Prepare("INSERT INTO distribution (id, show_id, title, slug) VALUES ($1, $2, $3, $4);"); err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	if _, err = stmt.Exec(distribution.Id, distribution.ShowId, distribution.Title, distribution.Slug); err != nil {
		return err
	}

	return nil
}

func (adapter *PostgresDistributionOutAdapter) updateDistributionEntry(distribution *model.Distribution, transaction *sql.Tx) (err error) {
	var stmt *sql.Stmt

	if stmt, err = transaction.Prepare("UPDATE distribution SET show_id = $2, title = $3, slug = $4 WHERE id = $1;"); err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	if _, err = stmt.Exec(distribution.Id, distribution.ShowId, distribution.Title, distribution.Slug); err != nil {
		return err
	}

	return nil
}

func (adapter *PostgresDistributionOutAdapter) ExistsByTitleOrSlug(title string, slug string) bool {
	query := "SELECT EXISTS(SELECT 1 FROM distribution where title = $1 or slug = $2)"
	row := adapter.db.QueryRow(query, title, slug)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (adapter *PostgresDistributionOutAdapter) GetDistributionOrNil(id string) (distribution *model.Distribution, err error) {
	query := "SELECT * FROM distribution where id = $1"
	row := adapter.db.QueryRow(query, id)

	distribution = &model.Distribution{}
	if err = row.Scan(&distribution.Id, &distribution.ShowId, &distribution.Title, &distribution.Slug); err != nil {
		return nil, nil
	}

	if err = adapter.parseEpisodes(distribution); err != nil {
		return nil, err
	}

	return distribution, nil
}

func (adapter *PostgresDistributionOutAdapter) parseEpisodes(distribution *model.Distribution) error {
	query := "SELECT episode_id FROM episodes_distributions where distribution_id = $1"
	rows, _ := adapter.db.Query(query, distribution.Id)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var episodeId string
		if err := rows.Scan(&episodeId); err != nil {
			return err
		}
		distribution.Episodes = append(distribution.Episodes, episodeId)
	}
	return nil
}

func (adapter *PostgresDistributionOutAdapter) UpdateDistribution(distribution *model.Distribution) error {
	transaction, _ := adapter.db.Begin()
	defer func(transaction *sql.Tx) {
		_ = transaction.Rollback()
	}(transaction)

	if err := adapter.updateDistributionEntry(distribution, transaction); err != nil {
		return err
	}
	if err := adapter.updateShowDistributionMappingEntry(distribution, transaction); err != nil {
		return err
	}
	if err := adapter.deleteEpisodeDistributionMappingEntries(distribution, transaction); err != nil {
		return err
	}
	if err := adapter.insertEpisodeDistributionMappingEntries(distribution, transaction); err != nil {
		return err
	}
	_ = transaction.Commit()
	return nil
}

func NewPostgresDistributionRepository(db *sql.DB) *PostgresDistributionOutAdapter {
	return &PostgresDistributionOutAdapter{db: db}
}
