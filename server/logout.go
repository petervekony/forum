package server

import (
	"fmt"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	//check if user is logged in
	// else nothing to log out
	// check if session is alive
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		// No session found, show sign up page
		//handle error
		fmt.Fprintln(w, err.Error())
	}
	// Check if user is logged in
	if uid == "0" {
		// User is not logged in, nothing to log out and can only see basic page without posts ability
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	}
	// if user is logged in, log out
	// delete the session
	cookie, err := sessionManager.deleteSession(w, r)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}

	// remove the cookie
	http.SetCookie(w, cookie)

	// redirect to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}
