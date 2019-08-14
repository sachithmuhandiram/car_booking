package notification

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/sendEmail", sendSignUpEmail)
}

func sendSignUpEmail(t_res http.ResponseWriter, t_req *http.Request) {
	log.Println("Hurray, API gateway works")
	//notification.SendNotification(email)
}
