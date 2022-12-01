package server

import (
	"fmt"
	logger "gritface/log"
	"net/http"
	"text/template"
)

// function handles the front page
func FrontPage(w http.ResponseWriter, r *http.Request) {

	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		logger.WTL(err.Error(), true)
	}

	if r.URL.Path == "/" { // TBC for session check
		if uid != "0" && r.Method == "POST" {
			// user is logged in redirect to front page with posts
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		tmpl, err := template.ParseFiles("server/public_html/index.html")
		if err != nil {
			logger.WTL(err.Error(), true)
			ErrorPage(w, 500)
		}
		dyn_content := make(map[string]string)
		dyn_content["name"] = uid
		tmpl.Execute(w, dyn_content)
	} else if r.URL.Path == "/posts" {
		posts, err := getPosts(r, uid, 0)
		if err != nil {
			logger.WTL(err.Error(), true)
			ErrorPage(w, 500)
		}
		w.Write([]byte(posts))
	} else if r.URL.Path == "/signup" {
		signupMsg, signupSuccess := SignUp(w, r)
		// Format response json
		writeMsg := fmt.Sprintf("{\"message\": \"%v\", \"status\": %v}", signupMsg, signupSuccess)
		w.Write([]byte(writeMsg))
	} else if r.URL.Path == "/login" {
		loginMsg, loginSuccess := Login(w, r)
		// Format response json
		writeMsg := fmt.Sprintf("{\"message\": \"%v\", \"status\": %v}", loginMsg, loginSuccess)
		w.Write([]byte(writeMsg))
	} else if r.URL.Path == "/loginSuccess" { // Active session front page
		if uid == "0" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		templ, err := template.ParseFiles("server/public_html/user.html")
		if err != nil {
			logger.WTL(err.Error(), true)
			ErrorPage(w, 500)
		}
		pic, err := GetProfilePic(uid)
		if err != nil {
			logger.WTL(err.Error(), true)
			ErrorPage(w, 500)
		}
		addPageInfo := map[string]string{
			"user_pic": pic,
		}
		templ.Execute(w, addPageInfo)
	} else if r.URL.Path == "/checkSession" {
		if err != nil {
			// No session found, show login page
			writeMsg := fmt.Sprintf("{\"status\": %v}", false)
			w.Write([]byte(writeMsg))
			return
		}
		// Check if user is logged in
		if uid != "0" {
			// User is logged in, redirect to front page
			writeMsg := fmt.Sprintf("{\"status\": %v}", true)
			w.Write([]byte(writeMsg))
		} else {
			writeMsg := fmt.Sprintf("{\"status\": %v}", false)
			w.Write([]byte(writeMsg))
		}
	} else if r.URL.Path == "/logout" {
		Logout(w, r)
	} else if r.URL.Path == "/getUser" {
		data, status := GetUserInfo(w, r)
		if status {
			w.Write([]byte(data))
		} else {
			logger.WTL(data, true)
			w.Write([]byte("{\"Username\": 0}"))
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else if r.URL.Path == "/addPost" {
		if r.Method != "POST" {
			w.WriteHeader(400)
			ErrorPage(w, 400)
		}
		message, status := addPostText(w, r)
		// Format response json
		writeMsg := fmt.Sprintf("{\"message\": \"%v\", \"status\": %v}", message, status)
		w.Write([]byte(writeMsg))
	} else if r.URL.Path == "/addComment" {
		if r.Method != "POST" {
			w.WriteHeader(400)
			ErrorPage(w, 400)
		}
		message, status := addComment(w, r)
		// Format response json
		writeMsg := fmt.Sprintf("{\"message\": \"%v\", \"status\": %v}", message, status)
		w.Write([]byte(writeMsg))
	} else if r.URL.Path == "/getCategories" {
		if r.Method != "POST" {
			w.WriteHeader(400)
			ErrorPage(w, 400)
		}
		message, status := GetCategories(w, r)
		fmt.Println(message, status)
		// writeMsg := fmt.Sprintf("{\"message\": \"%v\", \"status\": %v}", message, status)
		w.Write([]byte(message))
		// For users choice of filtering posts
	} else if r.URL.Path == "/filtered" {
		allMyPost, err := getPosts(r, uid, 0)
		if err != nil {
			logger.WTL(err.Error(), true)
			ErrorPage(w, 500)
		}
		w.Write([]byte(allMyPost))
	} else if r.URL.Path == "/add_reaction" {
		message, err := addReaction(w, r)
		if err != nil {
			logger.WTL(err.Error(), true)
			ErrorPage(w, 400)
		}
		w.Write([]byte(message))
	} else if r.URL.Path == "/loadPosts" {
		lastPostID := GetLastPostID(w, r)
		posts, err := getPosts(r, uid, lastPostID)
		if err != nil {
			logger.WTL(err.Error(), true)
		}
		w.Write([]byte(posts))
	} else {
		w.WriteHeader(404)
		ErrorPage(w, 404)
	}
}
