package main

import (
	"fmt"
	d "gritface/database"
	s "gritface/server"
	"log"
	"net/http"
	"time"
)

func main() {
	start := time.Now()
	// checking if the database is exists or not
	// then creating it
	_, err := d.DatabaseExist()
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

	// DISPLAY INSERTED RECORDS
	// d.QueryResultDisplay(forumdb)
	fmt.Println()
	fs := http.FileServer(http.Dir("server/public_html/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", s.FrontPage)
	fmt.Println("forum loaded in ", time.Since(start))
	fmt.Println()
	fmt.Println("Server is running on port 80...")
	go log.Fatalln(http.ListenAndServe("0.0.0.0:80", nil))
}
