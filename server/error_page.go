package server

import (
	"fmt"
	"io"
	"text/template"
)

type ErrorMessage struct {
	ErrorCode int
	Message   string
}

var messages = map[int]ErrorMessage{
	400: {
		ErrorCode: 400,
		Message:   "ERROR: Bad Request",
	},
	404: {
		ErrorCode: 404,
		Message:   "ERROR: Page Not Found",
	},
	500: {
		ErrorCode: 500,
		Message:   "ERROR: Internal Server Error",
	},
}

// ErrorPage function executes the error.html template with the given error code
func ErrorPage(wr io.Writer, error_code int) {
	templ, err := template.ParseFiles("server/public_html/error.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	templ.Execute(wr, messages[error_code])
}