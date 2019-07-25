package main

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFiles("ui/send_verification_email.html"))

}

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":8080", nil)

}

func foo(reswt http.ResponseWriter, req *http.Request) {
	// err := req.ParseForm()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// templates.ExecuteTemplate(reswt, "../ui/login.html", req.Form)

	//	temp, _ := templates.ParseFiles("../ui/login.html")

	tmpl.ExecuteTemplate(reswt, "send_verification_email.html", nil)

}
