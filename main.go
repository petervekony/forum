package main

import (
	"fmt"
	d "gritface/database"
	logger "gritface/log"
	s "gritface/server"
	"net/http"
)

var limiter = s.NewIPRateLimiter(50, 40)

func init() {
	go limiter.CleanUpVisitorMap()
}

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

	// create server struct
	ser := &http.Server{
		Addr:    ":443",
		Handler: limiter.LimitMiddleware(http.DefaultServeMux),
	}

	// start server
	logger.WTL("Server listening on '"+ser.Addr+"'", true)

	err = ser.ListenAndServe()
	if err != nil {
		logger.WTL(err.Error(), false)
	}

}
