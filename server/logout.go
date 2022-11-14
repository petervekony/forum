package server

import "net/http"

func Logout(w http.ResponseWriter, r *http.Request) {
	// delete the cookie
	cookie := http.Cookie{Name: "session", Value: "", Path: "/"}
	http.SetCookie(w, &cookie)

	// redirect to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}
