package database

import (
	"database/sql"
	"errors"
)

func UpdateCategoriesData(db *sql.DB, data map[string]string, category_id string) error {
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
