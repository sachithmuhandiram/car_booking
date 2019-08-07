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

type userLoginStruct struct {
	eventID int
	token   string
	used    bool
}

// database connection
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/car_booking_users")

	if err != nil {
		log.Println("Can not open database connection")
	}
	return db
}

// InitialToken insert to initial login page
func InitialToken(userLoginCookie string) {
	db := dbConn()
	defer db.Close()

	dt := time.Now()

	dateTime := dt.Format("2006-01-01 15:04:05.99999")

	tokenInsert, tpError := db.Prepare("insert into user_initial_login (token,created_date) values(?,?)")

	if tpError != nil {
		log.Println("Error prepairing initial token to table")
	}

	inRes, inErr := tokenInsert.Exec(userLoginCookie, dateTime)

	if inErr != nil {
		log.Println("Coulnd insert data to table", inRes)
		panic(inErr.Error())
	}
	log.Println("Initial Token inserted ", userLoginCookie)

	defer tokenInsert.Close()

}

// UserLoginData taking and sends for processing
func UserLoginData(loginResponse http.ResponseWriter, loginRequest *http.Request) {

	var userSessioneventID int

	if loginRequest.Method != "POST" {
		log.Panic("Form data is not Post")
		http.Redirect(loginResponse, loginRequest, "/", http.StatusSeeOther)
	}

	cookie, cookieError := loginRequest.Cookie("login-cookie") // returns cookie or an error

	// check incoming cookie with db

	if cookieError != nil {
		log.Fatal("Cookies dont match")
	} else {
		log.Println("Got cookie : ", cookie)
		userName := loginRequest.FormValue("username")
		password := loginRequest.FormValue("password")
		rememberMe := loginRequest.FormValue("remember_me")
		fmt.Println("Rember me  : ", rememberMe)

		userLogin, eventID := userLoginProcessing(userName, password, cookie.Value)

		if userLogin {
			/*
				Password matches, insert jwt and details to user_session table.
				Update initial loging table setting used=1 and next event id
			*/
			jwt, jwtErr := GenerateJWT(cookie.Value, 30) // for now its 30min session

			if jwtErr != nil {
				log.Println("Can not generate jwt token", jwtErr)
			}

			http.SetCookie(loginResponse, &http.Cookie{
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

			http.Redirect(loginResponse, loginRequest, "/home", http.StatusSeeOther)
		} else {
			// This is where I need to modify not to generate new token for login
			http.Redirect(loginResponse, loginRequest, "/", http.StatusSeeOther)
		}
	}

	// if user credintials are ok

}

// user data processing functions
func userLoginProcessing(userName string, password string, cookie string) (bool, int) {
	var loginUser users
	var userInitial userLoginStruct

	db := dbConn()

	defer db.Close()

	// login page defined token checking
	loginTokenCheck := db.QueryRow("SELECT event_id,used FROM user_initial_login WHERE token=? and used=0", cookie).Scan(&userInitial.eventID, &userInitial.used)

	if loginTokenCheck != nil {
		log.Println("user_initial_login table read faild") // posible system error or hacking attempt ?
		log.Println(loginTokenCheck)
		return false, 0
	}

	// update initial user details table
	initialUpdate, initErr := db.Prepare("update user_initial_login set used=1 where event_id=?")

	if initErr != nil {
		log.Println("Couldnt update initial user table")
		return false, 0 // we shouldnt compare password
	}

	_, updateErr := initialUpdate.Exec(userInitial.eventID)

	if updateErr != nil {
		log.Println("Couldnt execute initial update")

	}
	log.Printf("Initial table updated for event id %d : ", userInitial.eventID)
	// end login page token checking

	readError := db.QueryRow("SELECT id,password FROM car_booking_users WHERE username=?", userName).Scan(&loginUser.id, &loginUser.password)
	defer db.Close()
	if readError != nil {
		//http.Redirect(res, req, "/", 301)
		log.Println("data can not be taken")

	}

	comparePassword := bcrypt.CompareHashAndPassword([]byte(loginUser.password), []byte(password))

	// https://stackoverflow.com/questions/52121168/bcrypt-encryption-different-every-time-with-same-input

	if comparePassword != nil {
		/*
			Here I need to find a way to make sure that initial token is not get created each time wrong username password

			Also Need to implement a way to restrict accessing after 5 attempts
		*/
		log.Println("Wrong user name password")
		return false, 0
	} //else {

	log.Println("Hurray")
	return true, userInitial.eventID
	//}

}

// PasswordHashing exports hased password
func PasswordHashing(pass []byte) string {

	hashedPass, err := bcrypt.GenerateFromPassword(pass, 8)
	if err != nil {
		log.Println(err)
	}
	return string(hashedPass)
}

/*
	GenerateJWT will generate a JWT for a user seassion. Here I use initially generate token as an input.
	Also the time needs to keep token (valid_duration).
*/
func GenerateJWT(initialToken string, validDuration int) (string, error) {

	loginKey := []byte(initialToken)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(validDuration))

	jwtToken, jwtErr := token.SignedString(loginKey)

	if jwtErr != nil {
		log.Println("Error creating jwt Token : ", jwtErr)
		return "", jwtErr
	}

	return jwtToken, nil
}
