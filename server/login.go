package server

import (
	"encoding/json"
	"fmt"
	d "gritface/database"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the password string as a hash to be stored in the database
func passwordMatch(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// type User struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

func Login(w http.ResponseWriter, r *http.Request) (string, bool) {
	req, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return "Error: reading json log in request from user", false
	}

	// Unmarshal
	var user d.Users
	err = json.Unmarshal(req, &user)
	if err != nil {
		return "Error: unsuccessful in unmarshaling log in data from user", false
	}
	// If logged in, redirect to front page
	// If not logged in, show sign up page
	// check if session is alive
	uid, err := sessionManager.checkSession(w, r)
	fmt.Println("at check session uid is", uid)
	if err != nil {
		// No session found, show login page
		//handle error
		// fmt.Fprintln(w, err.Error())
		return err.Error(), false
	}
	// Check if user is logged in
	if uid != "0" {
		// User is logged in, redirect to front page
		// fmt.Fprintf(w, "User is logged in")
		// http.Redirect(w, r, "/", http.StatusSeeOther)
		return "User is logged in", false
	}

	// parse the form
	// r.ParseForm()

	// get the email and password
	email := EscapeString(user.Email)
	password := user.Password

	// check if the email and password are ascii
	if !IsAscii(email) || !IsAscii(password) {
		// fmt.Fprintf(w, "Invalid email or password")
		return "Error: invalid email or password", false
	}

	// retrieve user password from database
	db, err := d.DbConnect()
	if err != nil {
		return err.Error(), false
	}

	loginUser := make(map[string]string)
	loginUser["email"] = email
	users, err := d.GetUsers(db, loginUser)
	if err != nil {
		fmt.Println(err.Error())
	}
	
	fmt.Println(loginUser)
	// if no or more than 1 record found, return error
	if len(users) != 1 {
		return "Error: email or password is not found!", false
	}

	// check if the password match the one with database
	if !passwordMatch(password, users[0].Password) {
		// fmt.Fprintf(w, "Invalid email or password")
		return "Error: email or password is not found", false
	}

	// set the session
	// uID, err := strconv.Atoi(uid)
	// if err != nil {
	// 	// fmt.Fprintln(w, err.Error())
	// 	return err.Error(), false
	// }

	num := users[0].User_id

	err = sessionManager.setSessionUID(num, w, r)
	if err != nil {
		// fmt.Fprintln(w, err.Error())
		return err.Error(), false
	}
	return strconv.Itoa(num), true
	// redirect to the home page
	// http.Redirect(w, r, "/", http.StatusFound)
}
