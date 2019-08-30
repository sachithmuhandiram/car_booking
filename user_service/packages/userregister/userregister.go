package userregister

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

// GetEmail user data retrival functions
func GetEmail(registerResponse http.ResponseWriter, registerRequest *http.Request) {
	validEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`) // regex to validate email address

	if registerRequest.Method != "POST" {
		log.Panic("Form data is not Post")
		http.Redirect(registerResponse, registerRequest, "/", http.StatusSeeOther)
	}
	email := registerRequest.FormValue("email")

	if validEmail.MatchString(email) {
		fmt.Println("Valida email")
		fmt.Println("Email address : ", email)

		sendRegisterEmail(email)
	} else {
		fmt.Println("Wrong email")
	}

}

func sendRegisterEmail(email string) {
	log.Println("Send register email")

	sendMail, err := http.Get("http://notification:7070/sendmail")
	//Post("http://notification:7070/sendmail", "application", bytes.NewBuffer([]byte(email)))

	if err != nil {
		log.Println("Couldnt post send email")
	}

	log.Println("send mail sent : ", sendMail)

}
