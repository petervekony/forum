package main

import (
	"fmt"
	d "gritface/database"
	s "gritface/server"
	"net/http"
)

func main() {
	// check if db exist
	_, err := d.DatabaseExist()
	if err != nil {
		fmt.Println(err.Error())
	}
	// setup file server
	fs := http.FileServer(http.Dir("server/public_html/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// setup page handlers
	http.HandleFunc("/", s.FrontPage)
	fmt.Println()

	// start server
	fmt.Println("Server is running on port 80...")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
