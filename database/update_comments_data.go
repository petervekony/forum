package database

import (
	"database/sql"
)

// UpdateCommentsData updates the comment body of the comment specified by the comment_id.
// If the comment_id row doesn't exist or the body cannot be updated for any reason, it returns an error.
func UpdateCommentsData(db *sql.DB, body string, comment_id string) error {
	search := "SELECT * FROM comments WHERE comment_id=?"
	err := db.QueryRow(search, comment_id).Scan()
	if err == sql.ErrNoRows {
		return err
	}
	query := "UPDATE comments SET body=? WHERE comment_id=?"
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statement.Exec(body, comment_id)
	return err
}
