package database

import (
	"database/sql"
	"errors"
)

// UpdateCategoriesData updates the category with the correct category_id with data specified
// in a map[string]string. If category_id doesn't exist or there is something wrong with the statement,
// the function returns an error.
func UpdateCategoriesData(db *sql.DB, data map[string]string, category_id string) error {
	// Checking if category_id exists
	search := "SELECT * FROM categories WHERE category_id=?"
	err := db.QueryRow(search, category_id).Scan()
	if err == sql.ErrNoRows {
		return err
	}

	query := "UPDATE categories SET"
	counter := 0
	for key, value := range data {
		if key == "category_id" {
			return errors.New("ERROR: change of ID not allowed")
		}
		if counter > 0 {
			query += ","
		}
		query += " " + key + "='" + value + "'"
		counter++
	}
	query += " WHERE category_id=" + category_id
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	return err
}
