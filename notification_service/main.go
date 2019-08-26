package main

import (
	"log"
	"net/http"

	"./packages/notification"
)

func main() {
	http.HandleFunc("/sendmail", sendMail)
	http.ListenAndServe("0.0.0.0:7070", nil)
	log.Println("started notification service")
}

func sendMail(res http.ResponseWriter, req *http.Request) {
	log.Println("Request came to notification")
	notification.SendNotification(res, req)
}
