package server

import (
	"fmt"
	d "gritface/database"
	"net/http"
)

// test for how to get reaction from the database with an http method for the js
// this is build with adrian get_user_info.go as reference
func GetReactionInfo(w http.ResponseWriter, r *http.Request) (string, bool) {
	// session checking
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		return err.Error(), false
	}
	if uid == "0" {
		return "invalid session", false
	}
	// connect to database
	db, err := d.DbConnect()
	if err != nil {
		return err.Error(), false
	}
	// get the reaction info from the db as with key user_id and value as the uid
	user_reaction := make(map[string]string)
	user_reaction["user_id"] = uid
	users, err := d.GetReaction(db, user_reaction)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(users[0].Reaction)
	fmt.Println(users[0].User_id)
	fmt.Println(users[0].Post_id)
	fmt.Println(users[0].Comment_id)

	return users[0].Reaction, true
}
