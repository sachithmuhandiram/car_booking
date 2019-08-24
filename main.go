package main

import (
	"html/template"
	"log"
	"net/http"

	cars "./car_module/packages"
	"./user_service/packages/cookiecheck"
	"./user_service/packages/userlogin"
	"./user_service/packages/userregister"
	_ "github.com/go-sql-driver/mysql"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("user_service/ui/*.html"))

	//	temp, _ := template.ParseFiles("ui/login.html")
	//fmt.Println(*temp)
}

func main() {
	http.HandleFunc("/", userLoginView)
	http.HandleFunc("/login", userlogin.UserLoginData)
	http.HandleFunc("/register", getEmailView)
	http.HandleFunc("/register_user", userregister.GetEmail)
	http.HandleFunc("/home", checkCookie)
	http.HandleFunc("/testView", testView)
	http.HandleFunc("/test", cars.ProcessPhotos)
	http.HandleFunc("/showCar", cars.ShowCars)
	http.ListenAndServe(":8080", nil)

}

func testView(reswt http.ResponseWriter, req *http.Request) {
	OutputHTML(reswt, "car_module/ui/uploadPhoto.html", nil)
}

// internal functions
func checkCookie(res http.ResponseWriter, req *http.Request) {

	cookieSet := cookiecheck.CheckCookie(req)
	log.Println("check cookie returns : ", cookieSet)

	if cookieSet {
		userHomeView(res, req)
	} else {
		userLoginView(res, req)
	}
}

// html views
func userLoginView(reswt http.ResponseWriter, req *http.Request) {
	/*
		Here I am using a custom cookie and its used before user login. This is add purely for learning purposes.
	*/
	userLoginCookie := userlogin.PasswordHashing([]byte("test_value")) // this must change to random later
	userlogin.InitialToken(userLoginCookie)

	http.SetCookie(reswt, &http.Cookie{
		Name:  "login-cookie",
		Value: userLoginCookie,
		Path:  "/",
	})

	// insert initial token to table

	OutputHTML(reswt, "user_service/ui/login.html", nil)
}

func getEmailView(registerResponse http.ResponseWriter, registerRequest *http.Request) {
	OutputHTML(registerResponse, "user_service/ui/send_verification_email.html", nil)
}

func userHomeView(homeResponse http.ResponseWriter, homeRequest *http.Request) {
	OutputHTML(homeResponse, "user_service/ui/user_home.html", nil)
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

// test function
func testFunction(t_res http.ResponseWriter, t_req *http.Request) {
	//http.Redirect(t_res, t_req, "/sendEmail", 200)
	//notification.SendNotification("Hi hi")
	//OutputHTML(t_res, "user_service/ui/test.html", nil)

}
