package main

import (
	"log"
	"net/http"
)

func main() {
	//log.Println("started notification service")
	http.HandleFunc("/sendmail", sendMail)
	log.Println(http.ListenAndServe("0.0.0.0:7070", nil))
}

func sendMail(res http.ResponseWriter, req *http.Request) {
	log.Println("Request came to notification")
	//notification.SendNotification(res, req)
}
