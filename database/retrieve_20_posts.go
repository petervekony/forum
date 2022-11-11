package database

import (
	"database/sql"
	"encoding/json"
)

type JSONData struct {
	Post_id      int    `json:"post_id"`
	User_id      int    `json:"user_id"`
	Heading      string `json:"heading"`
	Body         string `json:"body"`
	Closed_user  int    `json:"closed_user"`
	Closed_admin int    `json:"closed_admin"`
	Closed_date  string `json:"closed_date"`
	Insert_time  string `json:"insert_time"`
	Update_time  string `json:"update_time"`
	Image        string `json:"image"`
}

func Retrieve20Posts(db *sql.DB) (string, error) {
	var structSlice []JSONData
	query := "SELECT * FROM posts LIMIT 20"
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		rD := &JSONData{}
		err = rows.Scan(&rD.Post_id, &rD.User_id, &rD.Heading, &rD.Body, &rD.Closed_user, &rD.Closed_admin, &rD.Closed_date, &rD.Insert_time, &rD.Update_time, &rD.Image)
		if err != nil {
			return "", err
		}
		structSlice = append(structSlice, *rD)
	}
	res, err := json.Marshal(structSlice)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
