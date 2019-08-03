package userlogin

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type users struct {
	id       int    //`json:"id"`
	username string //`json:"username"`
	email    string //`json:"email"`
	password string //`json:"password"`
}

type user_login_struct struct {
	event_id int
	token    string
	used     bool
}

// database connection
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/car_booking_users")

	if err != nil {
		log.Println("Can not open database connection")
	}
	return db
}

// initial token insert
func InitialToken(user_login_cookie string) {
	db := dbConn()
	defer db.Close()

	dt := time.Now()
	//Format MM-DD-YYYY
	date_time := dt.Format("2006-01-01 15:04:05.99999")
	log.Println("current date time", date_time)
	token_insert, tp_error := db.Prepare("insert into user_initial_login (token,created_date) values(?,?)")

	if tp_error != nil {
		log.Println("Error prepairing initial token to table")
	}

	in_res, in_err := token_insert.Exec(user_login_cookie, date_time)

	if in_err != nil {
		log.Println("Coulnd insert data to table", in_res)
		panic(in_err.Error())
	}
	log.Println("Initial Token inserted ", user_login_cookie)

	defer token_insert.Close()

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

		user_login := userLogin(user_name, password, cookie)

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
func userLogin(user_name string, password string, cookie *http.Cookie) bool {
	var login_user users
	var user_initial user_login_struct

	db, _ := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/car_booking_users")

	// login page defined token checking
	login_token_check := db.QueryRow("SELECT id,password FROM user_initial_login WHERE token=?", cookie).Scan(&user_initial.event_id, &user_initial.used)

	if login_token_check != nil {
		log.Println("user_initial_login table read faild") // posible system error or hacking attempt ?
		return false
	}

	// end login page token checking

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
