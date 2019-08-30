package main

import (
	"log"
	"net/http"

	"./packages/notification"
)

func main() {
	log.Println("started notification service")
	http.HandleFunc("/sendmail", sendMail)
	err := http.ListenAndServe("0.0.0.0:7070", nil)

	if err != nil {
		log.Println("couldnt start")
	}

	log.Println("Printing err value : ", err)

}

func sendMail(res http.ResponseWriter, req *http.Request) {
	log.Println("Request came to notification")
	notification.SendNotification(res, req)
}
