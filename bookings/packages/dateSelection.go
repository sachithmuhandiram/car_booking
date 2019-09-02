package booking

import (
	"net/http"
	"text/template"
)

// Shows date selection page.
func ShowDates(res http.ResponseWriter, req *http.Request) {
	OutputHTML(res, "bookings/ui/dateSelection.html", nil)
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
