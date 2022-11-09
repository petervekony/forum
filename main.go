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
	// DISPLAY INSERTED RECORDS
	// d.QueryResultDisplay(forumdb)
	fmt.Println()
	fs := http.FileServer(http.Dir("server/public_html/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", s.ServerHandler)
	fmt.Println("forum loaded in ", time.Since(start))
	fmt.Println()
	fmt.Println("Server is running on port 80...")
	go log.Fatalln(http.ListenAndServe("0.0.0.0:80", nil))
}
