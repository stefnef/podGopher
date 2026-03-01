package user

import (
	"database/sql"
	"podGopher/adapter/outbound/repository/HelperSql"
	"podGopher/core/domain/model"
)

type PostgresUserOutAdapter struct {
	db *sql.DB
}

func (adapter *PostgresUserOutAdapter) SaveUser(user *model.User) (err error) {
	transaction, _ := adapter.db.Begin()
	defer func(transaction *sql.Tx) {
		_ = transaction.Rollback()
	}(transaction)

	if err = adapter.createUserEntry(user, transaction); err != nil {
		return err
	}
	if err = adapter.createShowUserMappingEntry(user, transaction); err != nil {
		return err
	}
	_ = transaction.Commit()
	return nil
}

func prepareAndExecuteForUserId(statement string, userId string, transaction *sql.Tx) (err error) {
	var stmt *sql.Stmt

	if stmt, err = transaction.Prepare(statement); err != nil {
		return err
	}

	defer HelperSql.CloseStatement(stmt)

	if _, err = stmt.Exec(userId); err != nil {
		return err
	}
	return nil
}

func (adapter *PostgresUserOutAdapter) createShowUserMappingEntry(user *model.User, transaction *sql.Tx) (err error) {
	var stmtInsertIntoShowUsers, stmtInsertIntoUserRoles *sql.Stmt

	if stmtInsertIntoShowUsers, err = transaction.Prepare("INSERT INTO show_users (show_id, user_id) VALUES ($1, $2);"); err != nil {
		return err
	}
	defer HelperSql.CloseStatement(stmtInsertIntoShowUsers)

	if stmtInsertIntoUserRoles, err = transaction.Prepare("INSERT INTO user_roles (show_id, user_id, role) VALUES ($1, $2, $3);"); err != nil {
		return err
	}
	defer HelperSql.CloseStatement(stmtInsertIntoUserRoles)

	for _, showRole := range user.ShowRoles {
		if _, err = stmtInsertIntoShowUsers.Exec(showRole.ShowId, user.Id); err != nil {
			return err
		}
		if _, err = stmtInsertIntoUserRoles.Exec(showRole.ShowId, user.Id, showRole.Role.Name()); err != nil {
			return err
		}
	}

	return nil
}

func (adapter *PostgresUserOutAdapter) deleteShowUserMappingEntries(user *model.User, transaction *sql.Tx) (err error) {
	stmtDeleteShowUsers := "DELETE FROM show_users WHERE user_id = $1;"
	stmtDeleteUserRoles := "DELETE FROM user_roles WHERE user_id = $1;"

	if err = prepareAndExecuteForUserId(stmtDeleteShowUsers, user.Id, transaction); err != nil {
		return err
	}

	if err = prepareAndExecuteForUserId(stmtDeleteUserRoles, user.Id, transaction); err != nil {
		return err
	}

	return nil
}

func (adapter *PostgresUserOutAdapter) createUserEntry(user *model.User, transaction *sql.Tx) (err error) {
	var stmt *sql.Stmt

	if stmt, err = transaction.Prepare("INSERT INTO podcast_user (id, username) VALUES ($1, $2);"); err != nil {
		return err
	}
	defer HelperSql.CloseStatement(stmt)

	if _, err = stmt.Exec(user.Id, user.Username); err != nil {
		return err
	}

	return nil
}

func (adapter *PostgresUserOutAdapter) updateUserEntry(user *model.User, transaction *sql.Tx) (err error) {
	var stmt *sql.Stmt

	if stmt, err = transaction.Prepare("UPDATE podcast_user SET username = $2 WHERE id = $1;"); err != nil {
		return err
	}
	defer HelperSql.CloseStatement(stmt)

	if _, err = stmt.Exec(user.Id, user.Username); err != nil {
		return err
	}

	return nil
}

func (adapter *PostgresUserOutAdapter) ExistsByUsername(showId string, username string) bool {
	query := "SELECT EXISTS(SELECT 1 from show_users su LEFT JOIN podcast_user pu on pu.id = su.user_id where su.show_id = $1 and pu.username = $2)"
	row := adapter.db.QueryRow(query, showId, username)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (adapter *PostgresUserOutAdapter) GetUserOrNil(id string) (user *model.User, err error) {
	query := `
        SELECT pu.id, pu.username, u_roles.show_id, u_roles.role 
        FROM podcast_user pu 
        LEFT JOIN user_roles u_roles ON pu.id = u_roles.user_id 
        WHERE pu.id = $1;
    `
	rows, _ := adapter.db.Query(query, id)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	return parseUserWithRoles(rows)
}

func parseUserWithRoles(rows *sql.Rows) (user *model.User, err error) {
	userRolesSet := make(map[string]string)

	for rows.Next() {
		var (
			userId   string
			username string
			showId   sql.NullString
			role     sql.NullString
		)

		if err = rows.Scan(&userId, &username, &showId, &role); err != nil {
			return nil, err
		}

		if user == nil {
			user = &model.User{
				Id:       userId,
				Username: username,
			}
		}

		if showId.Valid && role.Valid {
			userRolesSet[showId.String] = role.String
		}
	}

	if user != nil {
		for showId, role := range userRolesSet {
			user.ShowRoles = append(user.ShowRoles, model.ShowRole{
				ShowId: showId,
				Role:   model.ValueToRole(role),
			})
		}
	}

	return user, nil
}

func (adapter *PostgresUserOutAdapter) UpdateUser(user *model.User) error {
	transaction, _ := adapter.db.Begin()
	defer func(transaction *sql.Tx) {
		_ = transaction.Rollback()
	}(transaction)

	if err := adapter.updateUserEntry(user, transaction); err != nil {
		return err
	}
	if err := adapter.deleteShowUserMappingEntries(user, transaction); err != nil {
		return err
	}
	if err := adapter.createShowUserMappingEntry(user, transaction); err != nil {
		return err
	}
	_ = transaction.Commit()
	return nil
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserOutAdapter {
	return &PostgresUserOutAdapter{db: db}
}
