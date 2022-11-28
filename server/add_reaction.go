package server

import (
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
	checkQuery := "SELECT * FROM reaction WHERE user_id = " + uid + " AND post_id = " + sPostID + " AND comment_id = " + sComID + " AND reaction_id =" + reactID
	fmt.Println(checkQuery)
	result, err := db.Query(checkQuery)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	var newUserID, newPostID, newCommentID, newReactionID int
	if result.Next() {
		err = result.Scan(&newUserID, &newPostID, &newCommentID, &newReactionID)
		fmt.Println("newreactionid and nreactid are ", newReactionID, nreactID, newCommentID, newPostID, newUserID)
		if newReactionID == nreactID {
			fmt.Println("3 here")
			_, err = d.DeleteReaction(db, uid, sPostID, sComID, reactID)
			if err != nil {
				return nil, err
			}
			reactID = "0"
			retData["reaction_id"] = reactID
		}
	} else {
		_, err = d.InsertReaction(db, nUID, postID, comID, reactID)
		if err != nil {
			return nil, err
		}
	}


	// sends data to the frontend from here

	query := "SELECT reaction_id, COUNT (reaction_id) AS rCount FROM reaction WHERE post_id = " + sPostID + " AND comment_id = " + sComID + " GROUP BY reaction_id"
	res, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(query)
	fmt.Println(reactID)
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
	fmt.Println(string(strLine))
	return strLine, nil
}
