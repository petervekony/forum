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
	u := d.Users{
		User_id:    1,
		Name:       "femi",
		Email:      "femi@gritlab.ax",
		Password:   "letsgohomenow",
		Deactive:   0,
		User_level: "admin",
	}
	fmt.Println(u.GetEmail())
	fmt.Println(u.GetName())
	fmt.Println(d.RetrieveUsers(db, 1))

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
