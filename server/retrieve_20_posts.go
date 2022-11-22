package server

import (
	"encoding/json"
	"gritface/database"
	"strconv"
)

type JSONData struct {
	Post_id      int                  `json:"post_id"`
	User_id      int                  `json:"user_id"`
	Heading      string               `json:"heading"`
	Body         string               `json:"body"`
	Closed_user  int                  `json:"closed_user"`
	Closed_admin int                  `json:"closed_admin"`
	Closed_date  string               `json:"closed_date"`
	Insert_time  string               `json:"insert_time"`
	Update_time  string               `json:"update_time"`
	Image        string               `json:"image"`
	Comments     map[int]JSONComments `json:"comments"`
}

type JSONComments struct {
	CommentID int    `json:"comment_id"`
	Post_id   int    `json:"post_id"`
	User_id   int    `json:"user_id"`
	Body      string `json:"body"`
}

func Retrieve20Posts() (string, error) {
	db, err := database.DbConnect()

	if err != nil {
		return "", err
	}

	structSlice := make(map[int]JSONData)
	query := "SELECT * FROM posts LIMIT 20"
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	nextQuery := ""
	for rows.Next() {
		rD := &JSONData{
			Comments: make(map[int]JSONComments),
		}
		err = rows.Scan(&rD.Post_id, &rD.User_id, &rD.Heading, &rD.Body, &rD.Closed_user, &rD.Closed_admin, &rD.Closed_date, &rD.Insert_time, &rD.Update_time, &rD.Image)
		if err != nil {
			return "", err
		}
		postId := &rD.Post_id
		structSlice[*postId] = *rD

		thisPostId := &rD.Post_id
		nextQuery += " OR post_id=" + strconv.Itoa(*thisPostId)
	}

	// Query comments
	query = "SELECT comment_id, post_id, user_id, body FROM comments WHERE " + nextQuery[4:]
	rows, err = db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		row := &JSONComments{}
		err = rows.Scan(&row.CommentID, &row.Post_id, &row.User_id, &row.Body)
		if err != nil {
			return "", err
		}
		thisPostId := &row.Post_id
		structSlice[*thisPostId].Comments[row.CommentID] = *row

	}
	res, err := json.Marshal(structSlice)
	if err != nil {
		return "", err
	}

	// fmt.Println(structSlice)
	// fmt.Println(string(res))
	return string(res), nil
}
