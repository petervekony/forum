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

	_, err = d.InsertReaction(db, nUID, postID, comID, reactID)
	if err != nil {
		fmt.Println("error inserting reaction: ", err.Error())
		return nil, err
	}

	query := "SELECT reaction_id, COUNT (reaction_id) AS rCount FROM reaction WHERE post_id = " + sPostID + " AND comment_id = " + sComID + " GROUP BY reaction_id"
	res, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Close()
	var recID string
	var count int
	for res.Next() {
		err = res.Scan(&recID, &count)
		if err != nil {
			fmt.Println(err)
		}
	}

	retData := make(map[string]interface{})

	retData["status"] = true
	retData["1"] = count
	retData["2"] = count
	retData["userReaction"] = reactID
	strLine, err := json.Marshal(retData)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(strLine))
	return strLine, nil
}
