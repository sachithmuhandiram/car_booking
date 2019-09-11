package cars

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
)

type CarPhotos struct {
	Id    string `json:id`
	Image []byte `json:image`
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/car_booking_cars")

	if err != nil {
		log.Println("Can not open database connection")
	}
	return db
}

func ProcessPhotos(pResponse http.ResponseWriter, uploadRequest *http.Request) {

	if uploadRequest.Method != "POST" {
		log.Panic("Form data is not Post")
	}
	firstPhoto, firstHeader, firstErr := uploadRequest.FormFile("imageOne")

	if firstErr != nil {
		log.Println("Error parsing the image file")

	}
	// Getting a random text to photo
	newName := randomName()

	// Create new file
	savePhoto, errCreate := os.Create("car_module/gallery/" + newName + firstHeader.Filename)

	log.Println("Save photo : ", *savePhoto)
	if errCreate != nil {
		log.Fatal(errCreate)
	}

	_, errCopy := io.Copy(savePhoto, firstPhoto)

	if errCopy != nil {
		log.Println("Error copying inphoto", errCopy)
	}

	defer savePhoto.Close()
	defer firstPhoto.Close()

	insertImage("car_module/gallery/" + newName + firstHeader.Filename)

}

func insertImage(imagePath string) {
	db := dbConn()

	defer db.Close()
	insertImage, imageErr := db.Prepare("insert into test(id,image) values(?,?)")

	/*
		Should extract image name and then print it in errors
	*/
	if imageErr != nil {
		log.Println("Prepairing Image failed", imagePath)
		//return false
	}

	_, insertErr := insertImage.Exec(3, imagePath)

	if insertErr != nil {
		log.Println("Couldnt insert image to table", imagePath)
		//return false
	}
	//return true
}

// ShowCars show a car from database
func ShowCars(carResponse http.ResponseWriter, carRequest *http.Request) {

	log.Println("--- Starting ShowCars ---")
	db := dbConn()
	defer db.Close()
	id := 5
	var carPhoto CarPhotos

	cars := db.QueryRow("SELECT image FROM test WHERE id=?", id).Scan(&carPhoto.Image)

	if cars != nil {
		log.Println("Couldnt get image from DB", cars) // posible system error or hacking attempt ?
	}

	log.Println("Image taken from DB : ", carPhoto.Image)
	log.Println("Data type : ", reflect.TypeOf(carPhoto.Image))

	defer log.Println("--- End ShowCars ---")
}

func randomName() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)

	return s
}
