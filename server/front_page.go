package server

import (
	"fmt"
	"net/http"
	"text/template"
)

func FrontPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fmt.Println("Trying to reach unknown path ", r.URL.Path)
		//error404(w)
		return
	}

	uid, err := sessionManager.checkSession(w, r)
	if err != nil {
		// Handle error for session check fail
	}

	tmpl, err := template.ParseFiles("server/public_html/index.html")
	if err != nil {
		//error404(w)
		return
	}

	dyn_content := make(map[string]string)
	dyn_content["name"] = uid

	tmpl.Execute(w, dyn_content)
}
