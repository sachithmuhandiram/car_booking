package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFiles("ui/send_verification_email.html"))
	temp, _ := template.ParseFiles("ui/send_verification_email.html")
	fmt.Println(*temp)
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

	//tmpl.ExecuteTemplate(reswt, "send_verification_email.html", nil)
	OutputHTML(reswt, "ui/send_verification_email.html", nil)
}

// test function
// output html
func OutputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
