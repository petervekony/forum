package database

import (
	"database/sql"
	"errors"
	"time"
)

// UpdatePostsData receives the database, a map[string]string and a post_id and updates post_id's columns
// specified in the map. If the update is not possible for any reason, the function returns an error.
// It also updates the update_time column based on the current time.
func UpdatePostsData(db *sql.DB, data map[string]string, post_id string) error {
	time := time.Now().Format("2006-01-02 15:04:05")
	query := "UPDATE posts SET"
	count := 0
	notAllowed := map[string]bool{
		"post_id":     true,
		"user_id":     true,
		"insert_time": true,
		"update_time": true,
	}
	for key, val := range data {
		if notAllowed[key] {
			return errors.New("ERROR: update of IDs and times not allowed")
		}
		if count > 0 {
			query += ","
		}
		query += " " + key + "='" + val + "'"
		count++
	}
	query += ", update_time='" + time + "' WHERE post_id=" + post_id
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	return err
}
