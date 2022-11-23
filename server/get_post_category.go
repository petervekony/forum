package server

import (
	"fmt"
	d "gritface/database"
	"net/http"
	"strconv"
)

func getPostCategories(w http.ResponseWriter, r *http.Request) {
	var post newPosts

	// get the post id from the request
	pid := strconv.Itoa(post.Post_id)
	// pid might be empty if there is no post_id in the url
	if pid == "" {
		return
	}

	// connect to the db
	db, err := d.DbConnect()
	if err != nil {
		fmt.Println(err.Error())
	}
	// get the post categories from the db as with key post_id and value as the post_id
	post_category := make(map[string]string)
	post_category["post_id"] = pid

	catergory, err := d.GetPostCategories(db, post_category)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(catergory[0].Category_id)
}

