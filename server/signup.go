package server

import (
	"fmt"
	d "gritface/database"
	"net/http"
	"regexp"
	"strings"
)

// function to check if email is valid
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

// function to check if inputs are valid
func isAscii(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c <= 33 || c >= 127 {
			return false
		}
	}
	return true
}

// prevents sql injection
func escapeString(value string) string {
	var sb strings.Builder
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch c {
		case '\\', 0, '\n', '\r', '\'', '"':
			sb.WriteByte('\\')
			sb.WriteByte(c)
		case '\032':
			sb.WriteByte('\\')
			sb.WriteByte('Z')
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

func SignUp(w http.ResponseWriter, r *http.Request) {

	// If logged in, redirect to front page
	// If not logged in, show sign up page
	// check if session is alive
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		// No session found, show sign up page
		//handle error
		fmt.Fprintln(w, err.Error())
	}
	// Check if user is logged in
	if uid != "0" {
		// User is logged in, redirect to front page
		fmt.Fprintf(w, "User is logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// User is not logged in, show sign up page
	fmt.Fprintf(w, "User is not logged in")
	// user trying to sign up
	if r.Method == "POST" {
		// get form values
		// check if email is valid
		// check if email is already registered
		// check if password is valid
		// check if password and password confirmation match
		// check if no field is empty
		// create user
		// redirect to front page
		// first check if string is not sql injection
		username := escapeString(r.FormValue("username"))
		email := escapeString(r.FormValue("email"))
		// no need to escape password because its hashed before being stored
		password := r.FormValue("password")
		confirmPwd := r.FormValue("Confirm Password")

		// check if email is valid
		if !isEmailValid(email) {
			fmt.Fprintf(w, "Email is not valid")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		// check if pass and confirm pass match
		if password != confirmPwd {
			fmt.Fprintf(w, "Passwords do not match")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		// make sure that no fields are empty or non ascii
		if !isAscii(username) || !isAscii(email) || !isAscii(password) || !isAscii(confirmPwd) {
			fmt.Fprintf(w, "Please fill in all fields")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		// user level will by 1 by default i.e registered user
		userLevel := 1
		// connect to database
		db, err := d.DbConnect()
		if err != nil {
			fmt.Println(err.Error())
		}
		// create user
		rowUpdated, err := d.InsertUsers(db, username, email, password, userLevel)
		fmt.Println(rowUpdated)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		} else {
			fmt.Fprintf(w, "User created")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
}
