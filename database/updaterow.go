package database

import (
	"database/sql"
)

func UpdateUsername(db *sql.DB, name string, newName string) (int, error) {
	updateQuery := `UPDATE users SET name = ? WHERE name = ?`
	statement, err := db.Prepare(updateQuery)
	if err != nil {
		return 0, err
	}
	val, err := statement.Exec(newName, name)
	if err != nil {
		return 0, err
	}
	insertId, _ := val.LastInsertId()
	return int(insertId), nil
}
