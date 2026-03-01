package HelperSql

import "database/sql"

func CloseStatement(stmt *sql.Stmt) {
	_ = stmt.Close()
}
