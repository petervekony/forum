package server

import (
	"encoding/json"
	"fmt"
	d "gritface/database"
	"net/http"
)

type Info struct {
	Image    string `json:"Image"`
	Username string `json:"Username"`
}

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

	user := make(map[string]string)
	user["user_id"] = uid
	users, err := d.GetUsers(db, user)
	if err != nil {
		fmt.Println(err.Error())
	}

	var info Info
	info.Image = "https://img.icons8.com/office/2x/circled-user-male-skin-type-4.png"
	info.Username = users[0].Name;

	jsonInfo, err := json.Marshal(info)
    if err != nil {
      return err.Error(), false
    }
  fmt.Println(string(jsonInfo))
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

	var categories []string
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

		categories = append(categories, category.Category_Name)
	}

	res, err := json.Marshal(categories)
	if err != nil {
		return err.Error(), false
	}

	return string(res), true
}
