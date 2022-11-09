package server

import (
	"fmt"
	"net/http"
	"text/template"
)

func ServerHandler(w http.ResponseWriter, r *http.Request) {
	// checking if the path is not correct and returning 400
	if r.URL.Path != "/" {
		fmt.Fprint(w, "404 page not found")
		return
	}
	// parsing the html file
	t, err := template.ParseFiles("server/public_html/index.html")
	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w, "500 - Interal Server Error")
		return
	}
	t.Execute(w, nil)
}
