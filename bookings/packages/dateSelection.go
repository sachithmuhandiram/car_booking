package booking

import (
	"log"
	"net/http"
	"text/template"
)

// ShowDates shows date selection page.
func ShowDates(res http.ResponseWriter, req *http.Request) {
	OutputHTML(res, "bookings/ui/dateSelection.html", nil)
}

// SearchDates gets data from form and parse
func SearchDates(res http.ResponseWriter, req *http.Request) {
	log.Println("Sesrch date function")
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
