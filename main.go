package main

import (
	"html/template"
	"net/http"

	"./packages/userlogin"
	"./packages/userregister"
	_ "github.com/go-sql-driver/mysql"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("ui/*.html"))

	//	temp, _ := template.ParseFiles("ui/login.html")
	//fmt.Println(*temp)
}

func main() {
	http.HandleFunc("/", userLoginView)
	http.HandleFunc("/login", userlogin.UserLoginData)
	http.HandleFunc("/register", getEmailView)
	http.HandleFunc("/register_user", userregister.GetEmail)
	http.ListenAndServe(":8080", nil)

}

// html views
func userLoginView(reswt http.ResponseWriter, req *http.Request) {
	user_login_cookie := userlogin.PasswordHashing([]byte("test_value")) // this must change to random later

	http.SetCookie(reswt, &http.Cookie{
		Name:  "login-cookie",
		Value: user_login_cookie,
	})
	OutputHTML(reswt, "ui/login.html", nil)
}

func getEmailView(register_response http.ResponseWriter, register_request *http.Request) {
	OutputHTML(register_response, "ui/send_verification_email.html", nil)
}

// output html view generic
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
