package server

import "net/http"

// limit session to 1 active session per user

func limitSession(w http.ResponseWriter, r *http.Request) (string, bool) {
	// get the session cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		return err.Error(), false
	}
	// get the session id from the cookie
	sessionID := cookie.Value
	// get the session from the session manager
	session, err := sessionManager.checkSession(w, r)
	if err != nil {
		return err.Error(), false
	}
	// get the user id from the session
	uid := session.Values["user_id"]
	// get the user from the database
	user := make(map[string]string)
	user["user_id"] = uid.(string)
	db, err := d.DbConnect()
	if err != nil {
		return err.Error(), false
	}
	users, err := d.GetUsers(db, user)
	if err != nil {
		return err.Error(), false
	}
	// check if the user has a session
	if users[0].Session_id != "" {
		// if the user has a session, delete the session
		err = sessionManager.Delete(users[0].Session_id)
		if err != nil {
			return err.Error(), false
		}
	}
	// update the user with the new session id
	user["session_id"] = sessionID
	err = d.UpdateUser(db, user)
	if err != nil {
		return err.Error(), false
	}
	return "", true
}
