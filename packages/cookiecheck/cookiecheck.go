package cookiecheck

import (
	"log"
	"net/http"
)

func CheckCookie(homeRequest *http.Request) bool {
	// This function checks whether user has a valid token, if so send to home page or send to login
	cookie, cookieError := homeRequest.Cookie("user-cookie")

	if cookieError != nil {
		return false
	}
	log.Println("User cookie there : ", cookie)
	return true
}
