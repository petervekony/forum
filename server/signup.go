package server

import (
	"fmt"
	d "gritface/database"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// to store the password in the database
var HashedPassword string

// function to check if email is valid
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)
	return emailRegex.MatchString(e)
}

// function to check if inputs are valid
func IsAscii(s string) bool {
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
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// prevents sql injection
func EscapeString(value string) string {
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
func SignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("User signing up. Name:%v, email:%v, password:%v, confirm password:%v\n", r.FormValue("username"), r.FormValue("email"), r.FormValue("password"), r.FormValue("Confirm Password"))
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
	// fmt.Fprintf(w, "User is not logged in")
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
		fmt.Println(r.FormValue("username"))
		name := EscapeString(r.FormValue("username"))
		fmt.Println(name)
		email := EscapeString(r.FormValue("email"))
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
		if !IsAscii(name) || !IsAscii(email) || !IsAscii(password) || !IsAscii(confirmPwd) {
			// for testing
			fmt.Printf("name: %v email: %v password: %v confrimPwd: %v\n", IsAscii(name), IsAscii(email), IsAscii(password), IsAscii(confirmPwd))
			fmt.Fprintln(w, "Internal server error", http.StatusInternalServerError)
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		// user level will by 1 by default i.e registered user
		userLevel := 1
		// hash password
		HashedPassword, err := HashPassword(password)
		if err != nil {
			// handle error
			fmt.Fprintln(w, "Internal server error", http.StatusInternalServerError)
		}
		// connect to database
		db, err := d.DbConnect()
		if err != nil {
			fmt.Println(err.Error())
		}
		// check if username is already taken
		var user map[string]string
		users, err := d.GetUsers(db, user)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(users)
		for _, u := range users {
			if u.Name == name {
				fmt.Fprintf(w, "Username is already taken")
				http.Redirect(w, r, "/signup", http.StatusSeeOther)
				return
			}
		}
		// create user
		rowUpdated, err := d.InsertUsers(db, name, email, HashedPassword, userLevel)
		fmt.Println(rowUpdated)
		if err != nil {
			fmt.Fprintln(w, err.Error())
			// needed to finalize endpoint
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			// fmt.Fprintf(w, "Signup successful")
			// needed to finalize endpoint
			fmt.Printf("Sign up successful for %v\n", r.FormValue("username"))
			// check session for pop up log in modal/give log in message in some way
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}
	}
}
