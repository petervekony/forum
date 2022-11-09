package main

import (
	"fmt"
	g "gritface/database"
)

func main() {
	// I delete the file to avoid duplicated records.
	forumdb, err := g.DatabaseExist()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println()
	// DISPLAY INSERTED RECORDS
	g.QueryResultDisplay(forumdb)
}
