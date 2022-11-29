package server

import (
	"encoding/json"
	"fmt"
	d "gritface/database"
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// function to check if email is valid
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)
	return emailRegex.MatchString(e)
}

// function to check if inputs are valid
func isAscii(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c <= 33 || c > 127 {
			return false
		}
	}
	return true
}

// hash password returned the password string as a hash to be stored in the database
// this is done for security reasons
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
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

// function to sign up a user
func signUp(w http.ResponseWriter, r *http.Request) (string, bool) {
	req, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return "Error: reading json sign up request from user", false
	}
	// Unmarshal
	var user NewUser
	err = json.Unmarshal(req, &user)
	if err != nil {
		return "Error: unsuccessful in unmarshaling data from user", false
	}
	// If logged in, redirect to front page
	// If not logged in, show sign up page
	// check if session is alive
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		// No session found, show sign up page
		return err.Error(), false
	}
	// Check if user is logged in
	if uid != "0" {
		errMsg := "Error: user is logged in"
		return errMsg, false
	}
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
		name := escapeString(user.Name)
		email := strings.ToLower(escapeString(user.Email))
		// no need to escape password because its hashed before being stored
		password := user.Password
		confirmPwd := user.ConfirmPwd
		if name == "" || email == "" || password == "" || confirmPwd == "" {
			errMsg := "Fill in all the required fields"
			return errMsg, false
		}
		// check if email is valid
		if !isEmailValid(email) {
			errMsg := "Email is not valid"
			return errMsg, false
		}
		// check if pass and confirm pass match
		if password != confirmPwd {
			errMsg := "Passwords do not match"
			return errMsg, false
		}
		// make sure that no fields are empty or non ascii
		if !isAscii(name) || !isAscii(email) || !isAscii(password) || !isAscii(confirmPwd) {
			// for testing
			errMsg := "Non-ascii characters found"
			return errMsg, false
		}
		// user level will by 1 by default i.e registered user
		userLevel := 1
		// hash password
		HashedPassword, err := hashPassword(password)
		if err != nil {
			errMsg := "Internal server error in hashing password"
			return errMsg, false
		}
		// connect to database
		db, err := d.DbConnect()
		if err != nil {
			return err.Error(), false
		}
		// check if username is already taken
		user := make(map[string]string)
		user["name"] = name
		fmt.Println(user)
		users, err := d.GetUsers(db, user)
		if err != nil {
			fmt.Println("???", err)
			// return err.Error(), false
		}
		for _, u := range users {
			if u.Name == name {
				errMsg := "Username is already taken"
				return errMsg, false
			}
		}
		// create user
		_, err = d.InsertUsers(db, name, email, HashedPassword, userLevel)
		if err != nil {
			if strings.Contains(err.Error(), "name") {
				return "Error: user name is already used.", false
			}
			if strings.Contains(err.Error(), "email") {
				return "Error: email is already used.", false
			}
		} else {
			fmt.Printf("Sign up successful for %v\n", name)
			successMsg := "User has been successfully registered"
			return successMsg, true
		}
	}
	return "Wrong method", false
}
