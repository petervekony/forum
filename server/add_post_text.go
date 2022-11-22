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

type newPosts struct {
	Body    string `json:"postBody"`
	Heading string `json:"postHeading"`
	Categories []string `json:"postCats"`
}

func addPostText(w http.ResponseWriter, r *http.Request) (string, bool) {
	req, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return "Error: reading json log in request from user", false
	}

	// Unmarshal
	var post newPosts
	err = json.Unmarshal(req, &post)
	if err != nil {
		return "Error: unsuccessful in unmarshaling log in data from user", false
	}
	fmt.Println(post)
	// If logged in, redirect to front page
	// If not logged in, show sign up page
	// check if session is alive
	uid, err := sessionManager.checkSession(w, r)
	fmt.Println("Adding post, check session uid is", uid)
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
	postID, err := d.InsertPost(db, uID, post.Heading, post.Body, time.Now().String()[0:19], "new post image")
	if err != nil {
		return err.Error(), false
	}
	for _, category := range post.Categories {
		catMap := make(map[string]string)
		catMap["category_name"] = category
		categoryArr, err := d.GetCategories(db, catMap)
		if err != nil {
			return err.Error(), false
		}
		_, err = d.InsertPostCategory(db, postID, categoryArr[0].Category_id)
		if err != nil {
			return err.Error(), false
		}
	}
	postID_str := strconv.Itoa(postID)
	return postID_str, true
}