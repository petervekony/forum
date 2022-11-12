package main

import (
	"fmt"
	d "gritface/database"
	s "gritface/server"
	"net/http"
	"time"
)

func main() {
	start := time.Now()
	// checking if the database is exists or not
	// then creating it
	db, err := d.DatabaseExist()
	if err != nil {
		fmt.Println(err.Error())
	}
	user := make(map[string]string)
	user["name"] = "%te%"
	user["email"] = "%@gritlab.ax%"
	fmt.Println("user map is", user)
	users, _ := d.GetUsers(db, user)
	fmt.Println("user info is", users)

	// DISPLAY INSERTED RECORDS
	// d.QueryResultDisplay(forumdb)
	fmt.Println()
	fs := http.FileServer(http.Dir("server/public_html/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", s.FrontPage)
	fmt.Println("forum loaded in ", time.Since(start))
	fmt.Println()
	fmt.Println("Server is running on port 80...")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
