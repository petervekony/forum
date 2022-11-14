package server

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the password string as a hash to be stored in the database
func passwordMatch(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(w http.ResponseWriter, r *http.Request) {
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

	// check if the username and password match
	if !passwordMatch(password, HashedPassword) {
		fmt.Fprintf(w, "Invalid username or password")
		return
	}

	// set the cookie
	cookie := http.Cookie{Name: "session", Value: username, Path: "/"}
	http.SetCookie(w, &cookie)

	// redirect to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}
