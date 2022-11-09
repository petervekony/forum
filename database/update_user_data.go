package database

import (
	"database/sql"
	"errors"
)

// UpdateUserData receives the database, a map[string]string and a user_id and updates user_id's columns
// specified in the map. If the update is not possible for any reason, the function returns an error.
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
