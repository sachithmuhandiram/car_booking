package cars

import (
	"database/sql"
	"log"
	"net/http"
)

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/car_booking_cars")

	if err != nil {
		log.Println("Can not open database connection")
	}
	return db
}

func ProcessPhoto(pResponse http.ResponseWriter, processRequest *http.Request) {
	if processRequest.Method != "POST" {
		log.Panic("Form data is not Post")
	}

	processRequest.ParseMultipartForm(10 << 20) // max 10MB
	photoOne, photoOneHandler, pOneError := processRequest.FormFile("imageOne")

	if pOneError != nil {
		log.Println("Error uploading photo")
	}
	// photoTwo := processRequest.FormValue("imageTwo")
	// photoThree := processRequest.FormValue("imageThree")
	// photoFour := processRequest.FormValue("imageFour")
	// photoFive := processRequest.FormValue("imageFive")

	//photos := []string{photoOne, photoTwo, photoThree, photoFour, photoFive}

	//for i := 0; i < 5; i++ {
	InsertCar(photoOneHandler.Filename)
	//}
	defer photoOne.Close()
}

func InsertCar(imagePath string) {
	db := dbConn()
	insertImage, imageErr := db.Prepare("insert into test(id,image) values(?,?)")

	/*
		Should extract image name and then print it in errors
	*/
	if imageErr != nil {
		log.Println("Prepairing Image failed", imagePath)
		//return false
	}

	_, insertErr := insertImage.Exec(1, imagePath)

	if insertErr != nil {
		log.Println("Couldnt insert image to table", imagePath)
		//return false
	}
	//return true
}
