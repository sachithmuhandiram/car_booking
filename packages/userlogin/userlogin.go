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

	cookie, cookie_error := login_request.Cookie("login-cookie")

	if cookie_error != nil {
		log.Fatal("Cookies dont match")
	} else {
		log.Println("Got cookie : ", cookie)
		user_name := login_request.FormValue("username")
		password := login_request.FormValue("password")
		remember_me := login_request.FormValue("remember_me")
		//p := []byte(password)
		//hashed_password := passwordHashing(p)

		fmt.Println("Rember me  : ", remember_me)

		userLogin(user_name, password)
	}
}

// user data processing functions
func userLogin(user_name string, password string) {
	var login_user users

	db, _ := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/car_booking_users")

	read_error := db.QueryRow("SELECT id,password FROM car_booking_users WHERE username=?", user_name).Scan(&login_user.id, &login_user.password)

	if read_error != nil {
		//http.Redirect(res, req, "/", 301)
		log.Println("data can not be taken")
		return
	}

	compare_password := bcrypt.CompareHashAndPassword([]byte(login_user.password), []byte(password))

	// https://stackoverflow.com/questions/52121168/bcrypt-encryption-different-every-time-with-same-input

	if compare_password != nil {
		log.Println("Wrong user name password")
	} else {
		log.Println("Hurray")
	}

	defer db.Close()

}

// internal functions
func PasswordHashing(pass []byte) string {

	hashed_pass, err := bcrypt.GenerateFromPassword(pass, 8)
	if err != nil {
		log.Println(err)
	}
	return string(hashed_pass)
}
