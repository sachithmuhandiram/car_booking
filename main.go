package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"./packages/userregister"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type users struct {
	id       int    `json:"id"`
	username string `json:"username"`
	email    string `json:"email"`
	password string `json:"password"`
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("ui/*.html"))

	//	temp, _ := template.ParseFiles("ui/login.html")
	//fmt.Println(*temp)
}

func main() {
	http.HandleFunc("/", userLoginView)
	http.HandleFunc("/login", userLoginData)
	http.HandleFunc("/register", getEmailView)
	http.HandleFunc("/register_user", userregister.GetEmail)
	http.ListenAndServe(":8080", nil)

}

// html views
func userLoginView(reswt http.ResponseWriter, req *http.Request) {
	OutputHTML(reswt, "ui/login.html", nil)
}

func getEmailView(register_response http.ResponseWriter, register_request *http.Request) {
	OutputHTML(register_response, "ui/send_verification_email.html", nil)
}

func userLoginData(login_response http.ResponseWriter, login_request *http.Request) {
	if login_request.Method != "POST" {
		log.Panic("Form data is not Post")
		http.Redirect(login_response, login_request, "/", http.StatusSeeOther)
	}

	user_name := login_request.FormValue("username")
	password := login_request.FormValue("password")
	remember_me := login_request.FormValue("remember_me")

	hashed_password := passwordHashing([]byte(password))
	// fmt.Println("User Name : ", user_name)
	// fmt.Println("Password : ", password)
	fmt.Println("Remmber me  : ", remember_me)
	// fmt.Println("hashed password   : ", hashed_password)

	userLogin(user_name, hashed_password)

}

// user data processing functions
func userLogin(user_name string, password string) {
	var login_user users

	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/car_booking_users")
	sqlQuery := `select id,password from car_booking_users where username=?;`
	if err != nil {
		log.Panic("Couldnt open database connection, userLogin")
	}

	val := db.QueryRow(sqlQuery, user_name)
	values := val.Scan(&login_user)

	fmt.Println("ID :", values) //, &login_user.password, &login_user.email))
	id := val.Scan(login_user.id)
	fmt.Println("Id : ", id)

	fmt.Println("User name and password", user_name, password)

	defer db.Close()

}

// internal functions
func passwordHashing(pass []byte) string {

	hashed_pass, err := bcrypt.GenerateFromPassword(pass, 8)
	if err != nil {
		log.Println(err)
	}

	return string(hashed_pass)
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
