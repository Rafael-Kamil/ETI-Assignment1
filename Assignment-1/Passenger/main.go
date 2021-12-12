package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-playground/validator/v10"

	"encoding/json"
	"fmt"
	"io/ioutil"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*This mathod to initialize the router*/
func initializeRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/passengers", GetPassengers).Methods("GET")
	router.HandleFunc("/passengers/Passgergettrips/{email}", PassengerGetTrips).Methods("GET")
	router.HandleFunc("/passengers/{email}", GetPassenger).Methods("GET")
	router.HandleFunc("/passengers", CreatePassenger).Methods("POST")
	router.HandleFunc("/passengers/{email}", UpdatePassenger).Methods("PUT")
	router.HandleFunc("/passengers/{email}", DeletePassenger).Methods("Delete")

	log.Fatal(http.ListenAndServe(":50001", router))

}

func main() {
	initialMigration()
	initializeRouter()

}

type Trip struct {
	TripID             int       `gorm:"primaryKey"`
	PassengerID        int       `json:"passengerid"`
	PickUpPoint        string    `json:"pickup" validate:"required"`
	DropoffPoint       string    `json:"dropoff" validate:"required"`
	DriverID           int       `json:"driverid"`
	StartTime          time.Time `json:"startDate"gorm:"autoCreateTime"`
	TripPassengeremail string    `json:"trip_passengeremail" validate:"required,email"`
	TripStatus         string    `json:"tripstatus"gorm:"default:waiting"`
	TripDriveremail    string    `json:"trip_driveremail" validate:"required,email"`
}
type Passenger struct {
	PassengerID int    `gorm:"primaryKey"`
	FirstName   string `json:"firstname" validate:"required"`
	Lastname    string `json:"lastname" validate:"required"`
	Phonenumber int    `json:"phonenumber" validate:"required" `
	Email       string `json:"email" validate:"required,email"`
}

var DB *gorm.DB
var err error

/*Here is the database connections string*/
const ADB = "root:00Nordic00@tcp(127.0.0.2:3306)/assignment1?charset=utf8mb4&parseTime=True&loc=Local"

/*Here is to to check if the connections to the database*/
func initialMigration() {
	DB, err = gorm.Open(mysql.Open(ADB), &gorm.Config{})

	if err != nil {
		fmt.Println(err.Error())
		panic("cant conenct to the Database Please check the coneections strings")
	}
	DB.AutoMigrate(&Passenger{})
}

//Get all registered Passenger
func GetPassengers(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var passenger []Passenger
	DB.Find(&passenger)
	json.NewEncoder(w).Encode(passenger)
}

//Here is a fucntions to get registerd Passenger by email
func GetPassenger(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var passenger Passenger
	err := DB.Where("email = ?", params["email"]).First(&passenger).Error
	if err != nil {
		//if user is not found
		fmt.Fprintf(w, "The email you enter is not registered")
		return
	} else {
		json.NewEncoder(w).Encode(passenger)
	}
}

//Here is a fucntions to create new Passenger
func CreatePassenger(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var passenger Passenger
	var dbpassenger Passenger
	json.NewDecoder(router.Body).Decode(&passenger)
	//to validate the inpute must be string not empty
	validate := validator.New()
	err2 := validate.Struct(passenger)
	if err2 != nil {
		fmt.Println(err2.Error())
		return
	}
	//Validate duplications of email address
	err := DB.Where("email = ?", passenger.Email).First(&dbpassenger).Error
	fmt.Println("Passenger: " + dbpassenger.FirstName)
	fmt.Println(err)
	if err == nil {
		fmt.Fprintf(w, "  The email you enter is already registerd  ")
		return
	}

	//Validate duplications of phone number
	err3 := DB.Where("phonenumber = ?", passenger.Phonenumber).First(&dbpassenger).Error
	fmt.Println("Passenger: " + dbpassenger.FirstName)
	if err3 == nil {
		fmt.Fprintf(w, "  The Phone number you enter is  already registerd    ")
		return
	}

	//if pass all validation Create passenger
	DB.Create(&passenger)
	json.NewEncoder(w).Encode(passenger)
}

//To be Advise
//Here is a fucntions to get registered passenger
func GetPassengerEmail(w http.ResponseWriter, router *http.Request) {
	var passenger Passenger
	var dbpassenger Passenger
	//Validate Registered  email address
	err := DB.Where("email = ?", passenger.Email).First(&dbpassenger).Error
	fmt.Println("Passenger: " + dbpassenger.FirstName)
	fmt.Println(err)
	if err != nil {
		DB.Create(&passenger)
		return
	} else {
		fmt.Printf("  The email you enter is already registerd  ")
	}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	DB.First(&passenger, params["email"])
	json.NewDecoder(router.Body).Decode(&passenger)
	DB.Save(&passenger)
	json.NewEncoder(w).Encode(passenger)

}

//Here is a fucntions to update passenger
func UpdatePassenger(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var passenger Passenger

	if err != nil {
		fmt.Printf("  The email you enter is not registerd  ")
		return
	} else {
		json.NewDecoder(router.Body).Decode(&passenger)
		DB.Model(&Passenger{}).Where("email=?", params["email"]).Updates(passenger)

		var newPassenger Passenger
		DB.Where("email=?", params["email"]).First(&newPassenger)
		json.NewEncoder(w).Encode(newPassenger)
		fmt.Printf("  Successfully update your account ")
	}
}

//Here is Fucntions to get Passenger trips
func PassengerGetTrips(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//var passenger Passenger
	//var trip Trip
	params := mux.Vars(router)
	email := params["email"]
	//err := DB.Where("email = ?", params["email"]).First(&passenger).Error
	if err != nil {
		//if user is not found
		fmt.Fprintf(w, "{The email you enter Is not registered}")
	} else {
		err, trips := GetPassengerTrips(email)
		fmt.Println(email)
		if err != nil {
			fmt.Print("no trips")
			return
		} else {
			//gettripsbyemail := DB.Find(&Passenger{}).Where("email = ?", trip.TripPassengeremail).Order("startDate")
			json.NewEncoder(w).Encode(trips)
		}
	}
}

//Here is the get trip url from trip.go
const TripUrl = "http://localhost:50003/trips"

//api call from trips.go to get passenger trips
func GetPassengerTrips(trip_passengeremail string) (error, []Trip) {
	url := TripUrl + "/" + trip_passengeremail
	var trips []Trip
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {

			errDecode := json.Unmarshal([]byte(body), &trips)
			fmt.Println()
			if errDecode != nil {
				return errDecode, trips
			} else {
				return nil, trips
			}
		} else {
			log.Fatal(err)
			return err, trips
		}
	} else {
		log.Fatal(err)
		return err, trips
	}
}

//Here is the fucntions to delete Passenger

func DeletePassenger(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var passenger Passenger

	err := DB.Where("email = ?", params["email"]).First(&passenger).Error

	if err != nil {
		fmt.Printf("  The email you enter is not registerd  ")

	} else {
		json.NewDecoder(router.Body).Decode(&passenger)
		json.NewEncoder(w).Encode("You are unable to Deleted the Account!")
	}
}
