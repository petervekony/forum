package server

import (
	"fmt"
	d "gritface/database"
	"net/http"
	"strconv"
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
	// get posts by user id and return the reactions to the posts
	users := make(map[string]string)
	users["user_id"] = uid

	// get post from user at uid
	user_posts, err := d.GetPosts(db, users)
	if err != nil {
		return err.Error(), false
	}
	// loop through the posts and get the reaction info
	for _, post := range user_posts {
		// get the reaction info from the db as with key post_id and value as the post_id
		post_reaction := make(map[string]string)
		// convert the post_id to string
		pid := strconv.Itoa(post.Post_id)
		// add the post_id to the map
		post_reaction["post_id"] = pid
		// get the reaction info from the db with the post_id
		reaction_info, err := d.GetReaction(db, post_reaction)
		if err != nil {
			return err.Error(), false
		}
		// print the reaction info since pid is unqiue
		fmt.Println(reaction_info[0].Reaction)
	}

	// get the reaction info from the db
	post_reaction := make(map[string]string)
	post_reaction["post_id"] = "1"
	// get the reaction info from the db

	posts, err := d.GetReaction(db, post_reaction)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(posts[0].Reaction)
	fmt.Println(posts[0].User_id)
	fmt.Println(posts[0].Post_id)
	fmt.Println(posts[0].Comment_id)

	return posts[0].Reaction, true
}

func GetPostReactions(w http.ResponseWriter, r *http.Request) (int, string, bool) {
	// session checking
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), false
	}
	var post Posts
	// connect to the db
	db, err := d.DbConnect()
	if err != nil {
		fmt.Println(err.Error())
	}
	// get the post_id from the url
	pid := strconv.Itoa(post.Post_id)
	// pid might be empty if there is no post_id in the url
	if pid == "" {
		return 0, "no post_id", false
	}
	// get the reaction info from the db as with key post_id and value as the post_id
	post_reaction := make(map[string]string)
	// convert the post_id to string

	// add the post_id to the map
	post_reaction["post_id"] = pid
	// get the reaction info from the db with the post_id
	reaction_info, err := d.GetReaction(db, post_reaction)
	if err != nil {
		return 0, err.Error(), false
	}
	// use this to prevent like and dislike from being added to the same post
	for _, reaction := range reaction_info {
		if strconv.Itoa(reaction.User_id) == uid {
			return 0, "already reacted", false
		}
	}
	// print the reaction info since pid is unqiue
	// len(reaction_info) should be the number of reaction for that post
	return len(reaction_info), reaction_info[0].Reaction, true
}
