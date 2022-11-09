package main

import (
	"fmt"
	g "gritface/database"
)

func main() {
	// checking if the database is exists or not
	// then creating it
	forumdb, err := g.DatabaseExist()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println()
	// DISPLAY INSERTED RECORDS
	g.QueryResultDisplay(forumdb)
}
