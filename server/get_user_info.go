package server

import (
	"encoding/json"
	"fmt"
	d "gritface/database"
	"net/http"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) (string, bool) {
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		return err.Error(), false
	}
	if uid == "0" {
		return "invalid session", false
	}

	db, err := d.DbConnect()
	if err != nil {
		return err.Error(), false
	}

	defer db.Close()

	user := make(map[string]string)
	user["user_id"] = uid
	users, err := d.GetUsers(db, user)
	if err != nil {
		fmt.Println(err.Error())
	}

	var info Info
	info.Image = "https://img.icons8.com/office/2x/circled-user-male-skin-type-4.png"
	info.Username = users[0].Name

	jsonInfo, err := json.Marshal(info)
	if err != nil {
		return err.Error(), false
	}

	return string(jsonInfo), true
}

type JSONCategories struct {
	Categories string `json:"categories"`
}

func GetCategories(w http.ResponseWriter, r *http.Request) (string, bool) {
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		return err.Error(), false
	}
	if uid == "0" {
		return "invalid session", false
	}

	db, err := d.DbConnect()
	if err != nil {
		return err.Error(), false
	}

	categories := make(map[int]string)
	rows, err := db.Query("select * from categories")
	if err != nil {
		return err.Error(), false
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var category d.Categories
		if err := rows.Scan(&category.Category_id, &category.Category_Name, &category.Closed); // Fetch the record
		err != nil {
			return err.Error(), false
		}

		categories[category.Category_id] = category.Category_Name
	}

	res, err := json.Marshal(categories)
	if err != nil {
		return err.Error(), false
	}

	return string(res), true
}

// func GetReactions(w http.ResponseWriter, r *http.Request) (string, bool) {
// 	// session checking
// 	uid, err := sessionManager.checkSession(w, r)
// 	if err != nil {
// 		return err.Error(), false
// 	}
// 	if uid == "0" {
// 		return "invalid session", false
// 	}

// 	db, err := d.DbConnect()
// 	if err != nil {
// 		return err.Error(), false
// 	}
// 	query := "select * from reactions WHERE user_id=" + uid
// 	row, err := db.Query(query)
// 	if err != nil {
// 		return err.Error(), false
// 	}
// 	defer row.Close()
// 	reactions := make(map[string]string)
// 	for row.Next() {
// 		var reaction d.Reaction
// 		if err := row.Scan(&reaction.User_id, &reaction.Post_id, &reaction.Comment_id, &reaction.Reaction_id); err != nil {
// 			return err.Error(), false
// 		}
// 		reactions[reaction.Reaction_id] = reaction.Reaction_id
// 	}
// 	res, err := json.Marshal(reactions)
// 	if err != nil {
// 		return err.Error(), false
// 	}
// 	return string(res), true

// }
