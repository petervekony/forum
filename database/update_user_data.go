package database

import (
	"database/sql"
	"errors"
)

func UpdateUserData(db *sql.DB, data map[string]string, user_id string) error {
	query := "UPDATE users SET"
	count := 0
	for key, val := range data {
		if key == "user_id" {
			return errors.New("ERROR: user_id update not allowed")
		}
		if count > 0 {
			query += ","
		}
		query += " " + key + "='" + val + "'"
		count++
	}
	query += " WHERE user_id=" + user_id
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	return err
}
