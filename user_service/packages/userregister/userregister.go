package userregister

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

// user data retrival functions
func GetEmail(register_response http.ResponseWriter, register_request *http.Request) {
	valid_email := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`) // regex to validate email address

	if register_request.Method != "POST" {
		log.Panic("Form data is not Post")
		http.Redirect(register_response, register_request, "/", http.StatusSeeOther)
	}
	email := register_request.FormValue("email")

	fmt.Println("Email address : ", email)

	if valid_email.MatchString(email) {
		fmt.Println("Valida email")
		sendRegisterEmail(email)
	} else {
		fmt.Println("Wrong email")
	}

}

func sendRegisterEmail(email string) {

}
