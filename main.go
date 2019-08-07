package main

import (
	"html/template"
	"log"
	"net/http"

	"./packages/cookiecheck"
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
	http.HandleFunc("/home", checkCookie)
	http.ListenAndServe(":8080", nil)

}

// internal functions
func checkCookie(res http.ResponseWriter, req *http.Request) {

	cookieSet := cookiecheck.CheckCookie(req)
	log.Println("check cookie returns : ", cookieSet)

	if cookieSet {
		userHomeView(res, req)
	} else {
		userLoginView(res, req)
	}
}

// html views
func userLoginView(reswt http.ResponseWriter, req *http.Request) {
	/*
		Here I am using a custom cookie and its used before user login. This is add purely for learning purposes.
	*/
	userLoginCookie := userlogin.PasswordHashing([]byte("test_value")) // this must change to random later
	userlogin.InitialToken(userLoginCookie)

	http.SetCookie(reswt, &http.Cookie{
		Name:  "login-cookie",
		Value: userLoginCookie,
		Path:  "/",
	})

	// insert initial token to table

	OutputHTML(reswt, "ui/login.html", nil)
}

func getEmailView(registerResponse http.ResponseWriter, registerRequest *http.Request) {
	OutputHTML(registerResponse, "ui/send_verification_email.html", nil)
}

func userHomeView(homeResponse http.ResponseWriter, homeRequest *http.Request) {
	OutputHTML(homeResponse, "ui/user_home.html", nil)
}

// OutputHTML view generic
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
