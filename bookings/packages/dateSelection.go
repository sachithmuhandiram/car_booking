package booking

import (
	"database/sql"
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

func dbConnCars() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/car_booking_cars")

	if err != nil {
		log.Println("Can not open database connection")
	}
	return db
}

// ShowDates shows date selection page.
func ShowDates(res http.ResponseWriter, req *http.Request) {
	OutputHTML(res, "bookings/ui/dateSelection.html", nil)
}

// AvailableVehicles gets data from form and parse
func AvailableVehicles(res http.ResponseWriter, req *http.Request) {
	db := dbConn()
	defer db.Close()

	log.Println("----- SearchDates started -----")
	if req.Method != "POST" {
		log.Panic("Search date data is not Post")
		http.Redirect(res, req, "/dateselection", http.StatusSeeOther)
	}

	/*
		Cookie matching to be done
	*/

	startDate := req.FormValue("startdate")
	endDate := req.FormValue("enddate")
	//vehicleType := req.FormValue("vtype")

	query := `SELECT vehicle_id FROM booking WHERE start_date > ? and start_date > ? 
				UNION 
			  SELECT vehicle_id FROM booking WHERE end_date < ?`

	bookingCheck, err := db.Query(query, startDate, endDate, startDate) //

	if err != nil {
		log.Println("Couldnt fetch available vehicles from booking table", bookingCheck)
	}

	var vehicleID int

	availableVehicles := []int{}

	for bookingCheck.Next() {
		err := bookingCheck.Scan(&vehicleID)

		if err != nil {
			log.Println("Scaning vehicle IDs failed")
		}

		availableVehicles = append(availableVehicles, vehicleID)
	}

	log.Println("Availble Vehicles : ", availableVehicles)

	showAvailableVehicles(res, availableVehicles)
	defer log.Println("---- End SearchDates ----")

}

// Internal functions
func showAvailableVehicles(r http.ResponseWriter, vID []int) {

	log.Println("--- Starting showAvailableVehicles ---")
	db := dbConnCars()
	defer db.Close()

	carPhoto := `SELECT image from test where id=?`
	for _, vid := range vID {
		//log.Println("Vehicle ID received : ", vid)
		avbCar, err := db.Query(carPhoto, vid)

		if err != nil {
			log.Println("Couldnt query car photo for vehicle_id", vid)
		}

		log.Println(avbCar)

	}
	defer log.Println("--- End showAvailableVehicles ---")
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
