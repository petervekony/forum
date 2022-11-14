package server

import (
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the password string as a hash to be stored in the database
func passwordMatch(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(w http.ResponseWriter, r *http.Request) {
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

	// parse the form
	r.ParseForm()

	// get the username and password
	username := EscapeString(r.Form.Get("username"))
	password := EscapeString(r.Form.Get("password"))

	// check if the username and password are valid
	if IsAscii(username) || IsAscii(password) {
		fmt.Fprintf(w, "Invalid username or password")
		return
	}

	// check if the password match the one with database
	if !passwordMatch(password, HashedPassword) {
		fmt.Fprintf(w, "Invalid username or password")
		return
	}

	// set the session
	uID, err := strconv.Atoi(uid)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	err = sessionManager.setSessionUID(uID, w, r)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}

	// redirect to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}
