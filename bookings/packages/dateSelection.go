package booking

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type booking struct {
	userID    int
	vehicleID int
	startDate time.Time
	endDate   time.Time
}

// database connection
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/car_booking")

	if err != nil {
		log.Println("Can not open database connection")
	}
	return db
}

// ShowDates shows date selection page.
func ShowDates(res http.ResponseWriter, req *http.Request) {
	OutputHTML(res, "bookings/ui/dateSelection.html", nil)
}

// SearchDates gets data from form and parse
func SearchDates(res http.ResponseWriter, req *http.Request) {
	db := dbConn()
	defer db.Close()

	var booking booking

	if req.Method != "POST" {
		log.Panic("Search date data is not Post")
		http.Redirect(res, req, "/dateselection", http.StatusSeeOther)
	}

	/*
		Cookie matching to be done
	*/

	startDate := req.FormValue("startdate")
	endDate := req.FormValue("enddate")
	vehicleType := req.FormValue("vtype")

	t := time.Now()
	today := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())

	bookingCheck := db.QueryRow("SELECT vehicle_id FROM booking WHERE start_date=? < ? and end_date< ?", startDate, today, endDate).Scan(&booking.vehicleID)

	if bookingCheck != nil {
		log.Println("Couldnt fetch available vehicles from booking table")
	}

	log.Println("Start date : ", startDate)
	log.Println("End date :", endDate)
	log.Println("Vehicle tyoe : ", vehicleType)

}

// OutputHTML view generic
func OutputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
