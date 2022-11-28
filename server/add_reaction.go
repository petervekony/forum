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
	fmt.Println("comID: ", comID)
	fmt.Println("reactID: ", reactID)
	fmt.Println("postID: ", postID)

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
	checkQuery := "SELECT reaction_id, COUNT (reaction_id) AS rCount FROM FROM reaction WHERE user_id = " + uid + " AND post_id = " + sPostID + " AND comment_id = " + sComID + " AND reaction_id ="+ reactID
	fmt.Println(checkQuery)
	_, err = db.Query(checkQuery)
	fmt.Println("error is", err)
	if err != sql.ErrNoRows {
		_, err = d.InsertReaction(db, nUID, postID, comID, reactID)
		if err != nil {
			fmt.Println("error inserting reaction: ", err.Error())
			return nil, err
		}
	} else {
		num, err := d.DeleteReaction(db, uid, sPostID, sComID, reactID)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("number of rows affected is", num)
		reactID = "0"
	}

	// sends data to the frontend from here

	query := "SELECT reaction_id, COUNT (reaction_id) AS rCount FROM reaction WHERE post_id = " + sPostID + " AND comment_id = " + sComID + " GROUP BY reaction_id"
	res, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(query)
	defer res.Close()
	var recID string
	var count int
	retData := make(map[string]interface{})
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
	fmt.Println(string(strLine))
	return strLine, nil
}
