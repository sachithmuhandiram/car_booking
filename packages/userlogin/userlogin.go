package userlogin

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type users struct {
	id       int    //`json:"id"`
	username string //`json:"username"`
	email    string //`json:"email"`
	password string //`json:"password"`
}

func UserLoginData(login_response http.ResponseWriter, login_request *http.Request) {
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
	fmt.Println("Rember me  : ", remember_me)
	// fmt.Println("hashed password   : ", hashed_password)

	userLogin(user_name, hashed_password)

}

// user data processing functions
func userLogin(user_name string, password string) {
	var login_user users

	db, _ := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/car_booking_users")

	// var databaseUsername string
	// var databasePassword string

	read_error := db.QueryRow("SELECT id,password FROM car_booking_users WHERE username=?", user_name).Scan(&login_user.id, &login_user.password)

	if read_error != nil {
		//http.Redirect(res, req, "/", 301)
		log.Println("data can not be taken")
		return
	}
	p := []byte(password)
	fmt.Println("database password : ", login_user.password)
	fmt.Println("my password", passwordHashing(p))
	//	hashed_password, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	compare_password := bcrypt.CompareHashAndPassword([]byte(login_user.password), p)
	// fmt.Println("Hashed user input", hashed_password)

	//fmt.Println("compared password :", compare_password)
	if compare_password != nil {
		log.Println("Wrong user name password")
	} else {
		log.Println("Hurray")
	}

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
