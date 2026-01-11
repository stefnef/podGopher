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
	return distribution, nil
}

func NewPostgresDistributionRepository(db *sql.DB) *PostgresDistributionOutAdapter {
	return &PostgresDistributionOutAdapter{db: db}
}
