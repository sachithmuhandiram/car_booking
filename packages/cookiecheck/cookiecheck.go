package cookiecheck

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

// UserSession table struct
type userSession struct {
	eventID int
	jwt     string
	expired int
}

// database connection
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/car_booking_users")

	if err != nil {
		log.Println("Can not open database connection")
	}
	return db
}

// CheckCookie fuction
func CheckCookie(homeRequest *http.Request) bool {

	var userCookies userSession
	// This function checks whether user has a valid token, if so send to home page or send to login
	cookie, cookieError := homeRequest.Cookie("user-cookie")

	if cookieError != nil {
		log.Println("No user-cookie set")
		return false
	}
	log.Println("User cookie there : ", cookie)
	// update user_session table
	db := dbConn()
	defer db.Close()

	//get data from user_session table
	loginTokenCheck := db.QueryRow("select event_id,expired from user_session where jwt=?", cookie.Value).Scan(&userCookies.eventID, &userCookies.expired)

	if loginTokenCheck != nil {
		log.Println("Couldnt read from user_session table")
	}

	if userCookies.expired == 1 {
		log.Println("User token expired")
		return false
	}
	dt := time.Now()

	dateTime := dt.Format("2006-01-01 15:04:05.99999")
	//update initial user details table
	sessionUpdate, sessionErr := db.Prepare("update user_session set last_login=? where event_id=?")

	if sessionErr != nil {
		log.Println("Couldnt update initial user table")
	}

	_, updateErr := sessionUpdate.Exec(dateTime, userCookies.eventID)

	if updateErr != nil {
		log.Println("Couldnt execute initial update")

	}
	log.Printf("Initial table updated for event id %d : ", userCookies.eventID)

	// log.Println(userSessioneventID)
	return true
}
