package userlogin

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

	dateTime := dt.Format("2006-01-01 15:04:05.99999")

	token_insert, tp_error := db.Prepare("insert into user_initial_login (token,created_date) values(?,?)")

	if tp_error != nil {
		log.Println("Error prepairing initial token to table")
	}

	in_res, in_err := token_insert.Exec(user_login_cookie, dateTime)

	if in_err != nil {
		log.Println("Coulnd insert data to table", in_res)
		panic(in_err.Error())
	}
	log.Println("Initial Token inserted ", user_login_cookie)

	defer token_insert.Close()

}

func UserLoginData(login_response http.ResponseWriter, login_request *http.Request) {

	var userSessioneventID int

	if login_request.Method != "POST" {
		log.Panic("Form data is not Post")
		http.Redirect(login_response, login_request, "/", http.StatusSeeOther)
	}

	cookie, cookie_error := login_request.Cookie("login-cookie") // returns cookie or an error

	// check incoming cookie with db

	if cookie_error != nil {
		log.Fatal("Cookies dont match")
	} else {
		log.Println("Got cookie : ", cookie)
		userName := login_request.FormValue("username")
		password := login_request.FormValue("password")
		remember_me := login_request.FormValue("remember_me")
		fmt.Println("Rember me  : ", remember_me)

		user_login, eventID := userLoginProcessing(userName, password, cookie.Value)

		if user_login {
			/*
				Password matches, insert jwt and details to user_session table.
				Update initial loging table setting used=1 and next event id
			*/
			jwt, jwt_err := GenerateJWT(cookie.Value, 30) // for now its 30min session

			if jwt_err != nil {
				log.Println("Can not generate jwt token", jwt_err)
			}

			http.SetCookie(login_response, &http.Cookie{
				Name:  "user-cookie",
				Value: jwt,
				Path:  "/home",
			})

			/*
				Inserting user_session and updating initial table
			*/
			t := time.Now()

			loginTime := t.Format("2006-01-01 15:04:05.99999")
			db := dbConn()

			defer db.Close()
			// inserting data to user_session
			insertSession, sessionErr := db.Prepare("insert into user_session (user_name,jwt,first_login,last_login) values(?,?,?,?)")

			if sessionErr != nil {
				log.Println("Couldnt insert data to user_session table")
			}

			_, insertErr := insertSession.Exec(userName, jwt, loginTime, loginTime)

			if insertErr != nil {
				log.Println("Couldnt execute insert to user_session table")

			}
			log.Printf("Data inserted to user_session table for User : %s : ", userName)

			getEventID, eventIDError := db.Prepare("select event_id from user_session where user_name=? and expired=0")

			if eventIDError != nil {
				log.Println("Couldnt get event_id from user_session table")
			}
			inErr := getEventID.QueryRow(userName).Scan(&userSessioneventID) // WHERE number = 1
			if inErr != nil {
				panic(inErr.Error()) // proper error handling instead of panic in your app
			}

			log.Println("Event ID in user session table :", userSessioneventID)

			//update initial user details table
			initialUpdate, initErr := db.Prepare("update user_initial_login set next_event_id=? where event_id=?")

			if initErr != nil {
				log.Println("Couldnt update initial user table")
			}

			_, updateErr := initialUpdate.Exec(userSessioneventID, eventID)

			if updateErr != nil {
				log.Println("Couldnt execute initial update")

			}
			log.Printf("Initial table updated for event id %d : ", eventID)

			log.Println(eventID)

			http.Redirect(login_response, login_request, "/home", http.StatusSeeOther)
		} else {
			// This is where I need to modify not to generate new token for login
			http.Redirect(login_response, login_request, "/", http.StatusSeeOther)
		}
	}

	// if user credintials are ok

}

// user data processing functions
func userLoginProcessing(user_name string, password string, cookie string) (bool, int) {
	var login_user users
	var user_initial user_login_struct

	db := dbConn()

	defer db.Close()

	// login page defined token checking
	login_token_check := db.QueryRow("SELECT event_id,used FROM user_initial_login WHERE token=? and used=0", cookie).Scan(&user_initial.event_id, &user_initial.used)

	if login_token_check != nil {
		log.Println("user_initial_login table read faild") // posible system error or hacking attempt ?
		log.Println(login_token_check)
		return false, 0
	}

	// update initial user details table
	initial_update, init_err := db.Prepare("update user_initial_login set used=1 where event_id=?")

	if init_err != nil {
		log.Println("Couldnt update initial user table")
		return false, 0 // we shouldnt compare password
	}

	_, update_err := initial_update.Exec(user_initial.event_id)

	if update_err != nil {
		log.Println("Couldnt execute initial update")

	}
	log.Printf("Initial table updated for event id %d : ", user_initial.event_id)
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
		/*
			Here I need to find a way to make sure that initial token is not get created each time wrong username password

			Also Need to implement a way to restrict accessing after 5 attempts
		*/
		log.Println("Wrong user name password")
		return false, 0
	} //else {

	log.Println("Hurray")
	return true, user_initial.event_id
	//}

}

// internal functions
func PasswordHashing(pass []byte) string {

	hashed_pass, err := bcrypt.GenerateFromPassword(pass, 8)
	if err != nil {
		log.Println(err)
	}
	return string(hashed_pass)
}

/*
	This will generate a JWT for a user seassion. Here I use initially generate token as an input.
	Also the time needs to keep token (valid_duration).
*/
func GenerateJWT(initial_token string, valid_duration int) (string, error) {

	login_key := []byte(initial_token)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(valid_duration))

	jwt_token, jwt_err := token.SignedString(login_key)

	if jwt_err != nil {
		log.Println("Error creating jwt Token : ", jwt_err)
		return "", jwt_err
	}

	return jwt_token, nil
}
