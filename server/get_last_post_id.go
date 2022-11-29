package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getLastPostID(w http.ResponseWriter, r *http.Request) int {
	req, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return 0
	}

	// Unmarshal
	var lastPost LastPost
	err = json.Unmarshal(req, &lastPost)
	if err != nil {
		return 0
	}

	return lastPost.ID
}
