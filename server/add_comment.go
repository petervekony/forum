package server

import (
	"encoding/json"
	"fmt"
	d "gritface/database"
	"io"
	"net/http"
	"strconv"
	"time"
)

func addComment(w http.ResponseWriter, r *http.Request) (string, bool) {
	req, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return "Error: reading json log in request from user", false
	}
	// Unmarshal
	var comment Comments
	err = json.Unmarshal(req, &comment)
	if err != nil {
		return "Error: unsuccessful in unmarshaling log in data from user", false
	}
	fmt.Println(comment)
	// If logged in, redirect to front page
	// If not logged in, show sign up page
	// check if session is alive
	uid, err := sessionManager.checkSession(w, r)
	fmt.Println("Adding comment, check session uid is", uid)
	if err != nil {
		// No session found, show login page
		//handle error
		// fmt.Fprintln(w, err.Error())
		return err.Error(), false
	}

	db, err := d.DbConnect()
	if err != nil {
		return err.Error(), false
	}
	uID, err := strconv.Atoi(uid)
	if err != nil {
		return err.Error(), false
	}
	if err != nil {
		return err.Error(), false
	}
	commentID, err := d.InsertComment(db, comment.Post_id, uID, comment.Body, time.Now().String())
	if err != nil {
		return err.Error(), false
	}
	commentID_str := strconv.Itoa(commentID)
	return commentID_str, true
}
