package main

import (
	"fmt"
	d "gritface/database"
	s "gritface/server"
	"net/http"
)

var limiter = s.NewIPRateLimiter(20, 10)

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

	// create server struct
	ser := &http.Server{
		Addr:    ":443",
		Handler: limiter.LimitMiddleware(http.DefaultServeMux),
	}
	// start server
	fmt.Println("Server is running on port 443...")

	// localhost.crt and localhost.key files were created using the following CLI commands:
	// openssl req  -new  -newkey rsa:2048  -nodes  -keyout localhost.key  -out localhost.csr
	// openssl  x509  -req  -days 365  -in localhost.csr  -signkey localhost.key  -out localhost.crt
	err = ser.ListenAndServeTLS("localhost.crt", "localhost.key")
	if err != nil {
		fmt.Println(err.Error())
	}

}
