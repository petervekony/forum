package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type SessionManager struct {
	sessions map[string]*SessionData
}
type SessionData struct {
	UId  string
	Misc map[string]string
}

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
			// Secure: true,
			// HttpOnly: true,
			MaxAge: 3600,
		}
		// Send cookie to client
		http.SetCookie(w, cookie)
		// Store cookie in session as uid=0 (unregistered user)
		var thisSessionData = &SessionData{
			UId: "0",
		}
		sm.sessions[ID] = thisSessionData
	}
	return sm.sessions[cookie.Value].UId, nil
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
func (sm *SessionManager) setSessionUID(uid int, w http.ResponseWriter, r *http.Request) error {
	thisSession, err := sm.isSessionSet(w, r)
	if err != nil {
		// Something wrong with cookie, return error
		return errors.New("could not retrieve cookie data")
	}
	// this is working right
	suid := strconv.Itoa(uid)
	// If uid is not 0, loop through current sessions, in case another session allready is logged in with uid, set it to zero
	if uid > 0 {
		for sessionId, setUid := range sm.sessions {
			if setUid.UId == suid && sessionId != thisSession.Value {
				fmt.Println("Unset previous session because of dual login")
				sm.sessions[sessionId].UId = "0"
			}
		}
	}
	sm.sessions[thisSession.Value].UId = strconv.Itoa(uid)
	return nil
}

// function to delete the session
func (sm *SessionManager) deleteSession(w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	// check if session is set
	cookie, err := sm.isSessionSet(w, r)
	if err != nil {
		return nil, err
	}
	// delete the session
	delete(sm.sessions, cookie.Value)
	// remove the cookie
	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		Secure: true,
		MaxAge: -1,
	}
	return cookie, nil
}
func init() {
	sessionManager = SessionManager{
		sessions: make(map[string]*SessionData),
	}
}
