package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

// function handles the front page
func FrontPage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling %v\n", r.URL.Path)
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		fmt.Println("error, session fucked up")
	}
	if r.URL.Path == "/" { // TBC for session check
		fmt.Println("cookies handling.")
		if uid != "0" && r.Method == "POST" {
			// user is logged in redirect to front page with posts
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		tmpl, err := template.ParseFiles("server/public_html/index.html")
		if err != nil {
			return
		}
		dyn_content := make(map[string]string)
		dyn_content["name"] = uid
		tmpl.Execute(w, dyn_content)
	} else if strings.Index(r.URL.Path, "/server/") == 0 {
		fmt.Printf("Handling %v\n", r.URL.Path[1:])
		script, err := os.ReadFile(r.URL.Path[1:])
		if err != nil {
			log.Fatal(err)
		}
		w.Write(script)
	} else if r.URL.Path == "/posts" {
		posts, err := Retrieve20Posts()
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(posts))
	} else if r.URL.Path == "/signup" {
		fmt.Printf("Signing up, path %v\n", r.URL.Path)
		signupMsg, signupSuccess := SignUp(w, r)
		if signupSuccess {
			fmt.Println(signupMsg)
			writeMsg := fmt.Sprintf("{\"message\": \"%v\", \"status\": %v}", signupMsg, true)
			w.Write([]byte(writeMsg))
		} else {
			// w.WriteHeader(400)
			fmt.Println(signupMsg)
			writeMsg := fmt.Sprintf("{\"message\": \"%v\", \"status\": %v}", signupMsg, false)
			w.Write([]byte(writeMsg))
			// w.Write([]byte(signupMsg))
		}
	} else if r.URL.Path == "/login" {
		fmt.Printf("Logging in, path %v\n", r.URL.Path)
		loginMsg, loginSuccess := Login(w, r)
		fmt.Println(loginMsg, loginSuccess)
		writeMsg := fmt.Sprintf("{\"message\": \"%v\", \"status\": %v}", loginMsg, loginSuccess)
		w.Write([]byte(writeMsg))
	} else if r.URL.Path == "/loginSuccess" {
		if uid == "0" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		fmt.Printf("User %v just logged in successfully.\n", r.FormValue("login_email"))
		templ, err := template.ParseFiles("server/public_html/user.html")
		if err != nil {
			log.Fatal(err)
		}
		pic, err := GetProfilePic(uid)
		if err != nil {
			log.Fatal(err)
		}
		addPageInfo := map[string]string{
			"user_pic": pic,
		}
		templ.Execute(w, addPageInfo)
		//w.Write(script)
	} else if r.URL.Path == "/checkSession" {
		if err != nil {
			// No session found, show login page
			//handle error
			// fmt.Fprintln(w, err.Error())
			writeMsg := fmt.Sprintf("{\"status\": %v}", false)
			w.Write([]byte(writeMsg))
		}
		// Check if user is logged in
		if uid != "0" {
			// User is logged in, redirect to front page
			// fmt.Fprintf(w, "User is logged in")
			// http.Redirect(w, r, "/", http.StatusSeeOther)
			writeMsg := fmt.Sprintf("{\"status\": %v}", true)
			w.Write([]byte(writeMsg))
		} else {
			writeMsg := fmt.Sprintf("{\"status\": %v}", false)
			w.Write([]byte(writeMsg))
		}
	} else if r.URL.Path == "/logout" {
		Logout(w, r)
	} else if r.URL.Path == "/getUser" {
		userInfo, status := GetUserInfo(w, r)
		if status {
			w.Write([]byte(userInfo))
		} else {
			fmt.Println(userInfo)
			w.Write([]byte("{\"Username\": 0}"))
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else if r.URL.Path == "/addPost" {
		message, status := addPostText(w, r)
		writeMsg := fmt.Sprintf("{\"message\": \"%v\", \"status\": %v}", message, status)
		w.Write([]byte(writeMsg))
	} else if r.URL.Path == "/addComment" {
		message, status := addComment(w, r)
		writeMsg := fmt.Sprintf("{\"message\": \"%v\", \"status\": %v}", message, status)
		w.Write([]byte(writeMsg))
	} else if r.URL.Path == "/getCategories" {
		message, status := GetCategories(w, r)
		fmt.Println(message, status)
		// writeMsg := fmt.Sprintf("{\"message\": \"%v\", \"status\": %v}", message, status)
		w.Write([]byte(message))
		// For users choice of filtering posts
	} else if r.URL.Path == "/filtered" {
		allMyPost, err := UserFilter(w, r, "userPosts", uid)
		if err != nil {
			fmt.Println(err)
		}
		w.Write([]byte(allMyPost))
	} else if r.URL.Path == "/add_reaction" {
		message, err := addReaction(w, r)
		if err != nil {
			fmt.Println(err)
		}
		w.Write([]byte(message))

	} else {
		w.WriteHeader(404)
		ErrorPage(w, 404)
	}
}
