package server

import (
	"encoding/json"
	"errors"
	"fmt"
	d "gritface/database"
	"net/http"
	"strconv"
)

func addReaction(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		return nil, err
	}
	if uid == "0" {
		return nil, errors.New("invalid session")
	}

	db, err := d.DbConnect()
	if err != nil {
		return nil, err
	}

	comID, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if err != nil {
		return nil, err
	}
	reactID, err := strconv.Atoi(r.URL.Query().Get("reaction_id"))
	if err != nil {
		return nil, err
	}
	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		return nil, err
	}
	fmt.Println("comID: ", comID)
	fmt.Println("reactID: ", reactID)
	fmt.Println("postID: ", postID)
	if postID == 0 || reactID < 1 || reactID > 2 {
		return nil, errors.New("invalid request")
	}
	_, err = d.InsertReaction(db, comID, reactID, postID, uid)
	if err != nil {
		return nil, err
	}
	//strLine := strconv.Itoa(line)

	retData := make(map[string]interface{})

	retData["status"] = true
	retData["change1"] = 1
	retData["change2"] = -1
	strLine, err := json.Marshal(retData)
	if err != nil {
		return nil, err
	}

	return strLine, nil
}
