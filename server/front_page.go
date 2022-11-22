package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func FrontPage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling %v\n", r.URL.Path)
	if r.URL.Path == "/" { // TBC for session check
		fmt.Println("cookies handling.")
		uid, err := sessionManager.checkSession(w, r)
		if err != nil {
			// Handle error for session check fail
		}

		if uid != "0" && r.Method == "GET" {
			// user is logged in redirect to front page with posts
			fmt.Fprintf(w, "User is logged in")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		tmpl, err := template.ParseFiles("server/public_html/index.html")
		if err != nil {
			//error404(w)
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
		fmt.Printf("User %v just logged in successfully.\n", r.FormValue("login_email"))
		// loginEmail := r.FormValue("login_email")
		script, err := os.ReadFile("server/public_html/user.html")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(script)
	} else if r.URL.Path == "/checkSession" {
		uid, err := sessionManager.checkSession(w, r)
		fmt.Println("at check session uid is", uid)
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
		Logout(w, r);
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
	} else {
		fmt.Println("Trying to reach unknown path ", r.URL.Path)
		// w.WriteHeader(404)
		// w.Write([]byte("404 Page not found."))
		// testing
		return
	}
}
