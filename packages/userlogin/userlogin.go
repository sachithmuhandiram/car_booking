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

	cookie, cookie_error := login_request.Cookie("login-cookie") // returns cookie or an error

	if cookie_error != nil {
		log.Fatal("Cookies dont match")
	} else {
		log.Println("Got cookie : ", cookie)
		user_name := login_request.FormValue("username")
		password := login_request.FormValue("password")
		remember_me := login_request.FormValue("remember_me")
		fmt.Println("Rember me  : ", remember_me)

		user_login := userLogin(user_name, password)

		fmt.Println("User Login", user_login)

		if user_login {
			http.Redirect(login_response, login_request, "/home", http.StatusSeeOther)
		} else {
			http.Redirect(login_response, login_request, "/", http.StatusSeeOther)
		}
	}

	// if user credintials are ok

}

// user data processing functions
func userLogin(user_name string, password string) bool {
	var login_user users

	db, _ := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/car_booking_users")

	read_error := db.QueryRow("SELECT id,password FROM car_booking_users WHERE username=?", user_name).Scan(&login_user.id, &login_user.password)
	defer db.Close()
	if read_error != nil {
		//http.Redirect(res, req, "/", 301)
		log.Println("data can not be taken")

	}

	compare_password := bcrypt.CompareHashAndPassword([]byte(login_user.password), []byte(password))

	// https://stackoverflow.com/questions/52121168/bcrypt-encryption-different-every-time-with-same-input

	if compare_password != nil {
		log.Println("Wrong user name password")
		return false
	} else {
		log.Println("Hurray")
		return true
	}

}

// internal functions
func PasswordHashing(pass []byte) string {

	hashed_pass, err := bcrypt.GenerateFromPassword(pass, 8)
	if err != nil {
		log.Println(err)
	}
	return string(hashed_pass)
}
