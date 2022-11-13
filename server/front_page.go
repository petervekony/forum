package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func FrontPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" { // TBC for session check
		fmt.Println("cookies handling.")
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
	} else if strings.Contains(r.URL.Path, "/server/") {
		script, err := os.ReadFile(r.URL.Path[1:])
		if err != nil {
			log.Fatal(err)
		}
		w.Write(script)
	} else if r.URL.Path == "/posts" {
		posts, err := Retrieve20Posts()
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(posts))
	} else {
		fmt.Println("Trying to reach unknown path ", r.URL.Path)
		//error404(w)
		return
	}
}
