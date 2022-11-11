package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type JSONData struct {
	post_id      int
	user_id      int
	heading      string
	body         string
	closed_user  int
	closed_admin int
	closed_date  string
	insert_time  string
	update_time  string
	image        string
}

// TODO: make it work, doesn't marshal at all ATM
func JSONThingy(db *sql.DB) (string, error) {
	var structSlice []JSONData
	query := "SELECT * FROM posts"
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		rD := &JSONData{}
		err = rows.Scan(&rD.post_id, &rD.user_id, &rD.heading, &rD.body, &rD.closed_user, &rD.closed_admin, &rD.closed_date, &rD.insert_time, &rD.update_time, &rD.image)
		if err != nil {
			return "", err
		}
		structSlice = append(structSlice, *rD)
	}
	fmt.Println(structSlice)
	res, err := json.Marshal(structSlice)
	if err != nil {
		return "", nil
	}
	fmt.Println("JSON: ", string(res))
	return string(res), nil
}
