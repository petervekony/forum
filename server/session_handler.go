package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type SessionManager struct {
	sessions map[string]string
}

// type SessionData struct {
// 	uid string
// 	ip  string
// }

var sessionManager SessionManager

// Check for valid session, if not create a new one. Return session user data
func (sm *SessionManager) checkSession(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := sm.isSessionSet(w, r)
	if err != nil {
		// Session is not set, create a new one

		ID := uuid.New().String() // Create a session id

		// Create session cookie
		cookie = &http.Cookie{
			Name:  "session",
			Value: ID,
			Path:  "/",
			// Secure:   true,
			// HttpOnly: true,
			MaxAge: 3600,
		}

		// Send cookie to client
		http.SetCookie(w, cookie)

		// Store cookie in session as uid=0 (unregistered user)
		sm.sessions[ID] = "0"
	}
	// here session value is 0 because user is not logged in
	fmt.Println("session value is", sm.sessions[cookie.Value])
	return sm.sessions[cookie.Value], nil
}

// This function check if a valid session is alive, returns session ID if alive
func (sm *SessionManager) isSessionSet(w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	c, err := r.Cookie("session")
	if err != nil {
		return nil, err
	}

	_, ok := sessionManager.sessions[c.Value] // Try to get session cookie value which will tell us if a valid session is open

	if ok {
		return c, nil
	} else {
		return nil, errors.New("cookie value not alive")
	}
}

// function set session UID
func (sm *SessionManager) setSessionUID(uid int, w http.ResponseWriter, r *http.Request) error {
	// user have logged in, set session UID
	thisSession, err := sm.isSessionSet(w, r)
	if err != nil {
		// Something wrong with cookie, return error
		return errors.New("could not retrieve cookie data")
	}
	// this is working right
	// here session value is 0
	sm.sessions[thisSession.Value] = strconv.Itoa(uid)
	fmt.Println("session value is now set to", sm.sessions[thisSession.Value])
	return nil
}

// functio delete session
func (sm *SessionManager) deleteSession(w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	// check if user is logged in
	cookie, err := sm.isSessionSet(w, r)
	if err != nil {
		return nil, err
	}

	// delete the session
	delete(sessionManager.sessions, cookie.Value)

	// remove the cookie
	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	return cookie, nil
}

func init() {
	sessionManager = SessionManager{
		sessions: make(map[string]string),
	}
}
