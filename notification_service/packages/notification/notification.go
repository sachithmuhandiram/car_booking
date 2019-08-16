package notification

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type emailDetails struct {
	emailData string `json: "emailDetails"`
	// fromEmail   string `json:"from"`
	// emailParser string `json:"parse"`
	// toEmail     string `json:"toemail"`
}

// This function reads the json file and pass values to SendNotification
func getCredintials() (string, string, string) {
	jsonFile, err := os.Open("notification_service/packages/notification/emailData.json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var email emailDetails
	json.Unmarshal(byteValue, &email)
	log.Println("User Type: " + email.emailDetails.fromEmail)

	return email.fromEmail, email.emailParser, email.toEmail

}

func SendNotification(res http.ResponseWriter, req *http.Request) {
	body := "hi hi"
	from, pass, to := getCredintials()

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent email")
}
