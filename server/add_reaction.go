package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	d "gritface/database"
	"net/http"
	"strconv"
	"strings"
)

func addReaction(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var recID string
	var count int
	retData := make(map[string]interface{})
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

	defer db.Close()

	sComID := r.URL.Query().Get("comment_id")
	comID, err := strconv.Atoi(sComID)
	if err != nil {
		return nil, err
	}
	reactID := r.URL.Query().Get("reaction_id")
	sPostID := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(sPostID)
	if err != nil {
		return nil, err
	}
	nreactID, err := strconv.Atoi(reactID)
	if err != nil {
		return nil, err
	}
	if postID == 0 || nreactID < 1 || nreactID > 2 {
		return nil, errors.New("invalid request")
	}
	nUID, err := strconv.Atoi(strings.TrimSpace(uid))
	if err != nil {
		return nil, err
	}
	if nUID == 0 {
		return nil, errors.New("invalid session")
	}

	queryStr := "SELECT count(post_id) as rowCount FROM reaction WHERE user_id=? AND post_id = ? AND comment_id = ? AND reaction_id = ?"
	queryChk, err := db.Prepare(queryStr)
	if err != nil {
		return nil, err
	}

	defer queryChk.Close()

	var rowCount string
	err = queryChk.QueryRow(nUID, postID, comID, reactID).Scan(&rowCount)

	if (err == sql.ErrNoRows || err == nil) && rowCount == "0" {
		// If no rows where found
		_, err = d.InsertReaction(db, nUID, postID, comID, reactID)
		if err != nil {
			return nil, err
		}
	} else if rowCount == "1" {
		// There is exactly a line like this allready, delete it
		_, err = d.DeleteReaction(db, uid, sPostID, sComID, reactID)
		if err != nil {
			return nil, err
		}
		reactID = "0"
	} else { // Something went wrong
		return nil, err
	}

	// sends data to the frontend from here

	query := "SELECT reaction_id, COUNT (reaction_id) AS rCount FROM reaction WHERE post_id = " + sPostID + " AND comment_id = " + sComID + " GROUP BY reaction_id"
	res, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Close()

	retData["rb1"] = 0
	retData["rb2"] = 0
	for res.Next() {
		err = res.Scan(&recID, &count)
		if err != nil {
			fmt.Println(err)
		}
		retData["rb"+recID] = count
	}

	retData["status"] = true
	retData["userReaction"] = reactID
	strLine, err := json.Marshal(retData)
	if err != nil {
		return nil, err
	}
	return strLine, nil
}
