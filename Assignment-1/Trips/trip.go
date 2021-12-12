package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-playground/validator/v10"

	"encoding/json"
	"fmt"
	"io/ioutil"

	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const baseURL = "http://localhost:50001/passengers"
const baseURL2 = "http://localhost:50002/drivers/available"

/*This mathod to initialize the router*/
func initializeRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/trips/{trip_passengeremail}", GetTripsByPassengerEmail).Methods("GET")
	router.HandleFunc("/trips/driver/{trip_driveremail}", GetTripsByDriverEamil).Methods("GET")
	//router.HandleFunc("/trip/{email}", Getrips).Methods("GET")
	router.HandleFunc("/trips/{email}", CreateTrip).Methods("POST")
	router.HandleFunc("/trips/processing/{trip_driveremail}", UpdateTripsStatusToProcessing).Methods("PUT")
	router.HandleFunc("/trips/completed/{trip_driveremail}", UpdateTripsStatusToComplete).Methods("PUT")

	log.Fatal(http.ListenAndServe(":50003", router))

}

//Here is to initialise Migrations and routers
func main() {
	initialMigration()
	initializeRouter()

}

//Trips struct
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

//Passenger struct
type Passenger struct {
	PassengerID int    `gorm:"primaryKey"`
	FirstName   string `json:"firstname" validate:"required"`
	Lastname    string `json:"lastname" validate:"required"`
	Phonenumber int    `json:"phonenumber" validate:"required" `
	Email       string `json:"email" validate:"required,email"`
}

//Driver Struct
type Driver struct {
	DriverID    int    `gorm:"primaryKey"`
	FirstName   string `json:"firstname" validate:"required"`
	Lastname    string `json:"lastname" validate:"required"`
	Phonenumber int    `json:"phonenumber" validate:"required" `
	Email       string `json:"email" validate:"required,email"`
	IcNum       string `json:"icnum" validate:"required" `
	LicenseNum  string `json:"licensenum" validate:"required" `
	Available   bool   `gorm:"type:bool; json:"available" validate:"required"`
}

var DB *gorm.DB
var err error

//Database coonnections string
const ADB = "root:00Nordic00@tcp(127.0.0.2:3306)/assignment1?charset=utf8mb4&parseTime=True&loc=Local"

func initialMigration() {
	DB, err = gorm.Open(mysql.Open(ADB), &gorm.Config{})
	/*Here is to to check if the connections to the database*/

	if err != nil {
		fmt.Println(err.Error())
		panic("cant conenct to the Database Please check the coneections strings")
	}
	DB.AutoMigrate(&Trip{})
}

//Here is a fucntions to update Trip status to completed. This will be use by Driver end trips
func UpdateTripsStatusToComplete(w http.ResponseWriter, router *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var trip Trip

	err := DB.Where("trip_driveremail = ?", params["trip_driveremail"]).Order("start_time desc").First(&trip).Error
	if err != nil {
		fmt.Printf("  The email you enter is not registerd  ")

	} else {
		json.NewDecoder(router.Body).Decode(&trip)
		DB.Model(&Trip{}).Where("trip_driveremail = ?", params["trip_driveremail"]).Update("trip_status", "completed")
		//DB.Save(&trip)
		json.NewEncoder(w).Encode(trip)
		fmt.Printf("  Successfully update your account ")

	}
}

//update trips status to Process. this will be used by Driver Starttrips fucntions
func UpdateTripsStatusToProcessing(w http.ResponseWriter, router *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var trip Trip
	//Here is to validate if driver have a trips
	err := DB.Where("trip_driveremail = ?", params["trip_driveremail"]).Order("start_time desc").First(&trip).Error
	if err != nil {
		fmt.Printf("  The email you enter is not registerd  ")

	} else {
		//If driver has a trips, change trips status to process
		json.NewDecoder(router.Body).Decode(&trip)
		DB.Model(&Trip{}).Where("trip_driveremail = ?", params["trip_driveremail"]).Update("trip_status", "process")
		json.NewEncoder(w).Encode(trip)
		fmt.Printf("  Successfully update your account ")

	}
}

//Get trips all trips by passenger email
func GetTripsByPassengerEmail(w http.ResponseWriter, router *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var trips []Trip
	//To validate if Passenger has trips
	err := DB.Where("trip_passengeremail = ?", params["trip_passengeremail"]).Order("start_time desc").Find(&trips).Error
	if err != nil {
		//if user is not found
		fmt.Fprintf(w, "The emaail you entered does not have trip")
	} else {
		//Print Passenger trips
		json.NewEncoder(w).Encode(trips)
	}
}

//Get trips by driver email
func GetTripsByDriverEamil(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var trips []Trip
	err := DB.Where("trip_driveremail = ?", params["trip_driveremail"]).Order("start_time desc").Find(&trips).Error
	fmt.Println(params["trip_driveremail"])
	if err != nil {
		fmt.Fprintf(w, "The emaail you entered does not have trip")
	} else {
		json.NewEncoder(w).Encode(trips)
	}
}

//create trips, and update driver avilable status.
func CreateTrip(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var trip Trip
	json.NewDecoder(router.Body).Decode(&trip)

	//to validate the inpute must be string not empty
	validate := validator.New()
	err2 := validate.Struct(trip)
	if err != nil {
		fmt.Println(err2.Error())
		return
	} else {
		fmt.Printf("The passenger is successfully Created")
	}

	params := mux.Vars(router)
	email := params["email"]
	//Validate registed user of email address
	err, passenger := GetPassengerIdByEmail(email)
	if err != nil {
		fmt.Fprintf(w, " The email you registerd is not registerd as passenger ")
	} else {
		//get the passenger ID and emails
		tripsemail := passenger.Email
		trip.TripPassengeremail = tripsemail
		passengerId := passenger.PassengerID
		trip.PassengerID = passengerId

		//get available driver id
		driverErr, driver := GetRandomDriverId()

		if driverErr != nil {
			fmt.Printf("No available driver")
		} else {
			//Get driver email and ID and store in the trip DriverId and DriverEamail
			driveremail := driver.Email
			trip.TripDriveremail = driveremail
			driverid := driver.DriverID
			trip.DriverID = driverid
			//update the driver status.
			//call the driver api update with the querry above
			updatedriverErr2 := UpdateDriverAvilablityToFalse(driveremail)
			if updatedriverErr2 != nil {
				fmt.Print("Error While updating Availability")
				return
			} else {
				DB.Create(&trip)
				json.NewEncoder(w).Encode(trip)
				//driveremail := driver.Email
				fmt.Printf("successfuly update driver")
				return
			}

		}
	}
}

//Update driver avilability to False
const UpdateDriverAvailToFalse = "http://localhost:50002/driver/available"

//Get Register passenger by UpdateDriver UpdateDriverAvailToFalse email
func UpdateDriverAvilablityToFalse(email string) error {
	request, err := http.NewRequest(http.MethodPut,
		UpdateDriverAvailToFalse+"/"+email, nil)

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return err
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
		return nil
	}
}

//API call to get Driver by email
func GetPassengerIdByEmail(email string) (error, Passenger) {
	url := baseURL + "/" + email
	var passenger Passenger
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {

			errDecode := json.Unmarshal([]byte(body), &passenger)

			if errDecode != nil {
				return errDecode, passenger
			} else {
				return nil, passenger
			}
		} else {
			log.Fatal(err)
			return err, passenger
		}
	} else {
		log.Fatal(err)
		return err, passenger
	}
}

//API call to get random Available driver ,  mean avilability = True
func GetRandomDriverId() (error, Driver) {
	//call driver API here
	url2 := baseURL2 + "/true"
	var driver Driver
	if resp, err := http.Get(url2); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			errDecode := json.Unmarshal([]byte(body), &driver)

			if errDecode != nil {
				return errDecode, driver
			} else {
				return nil, driver
			}
		} else {
			log.Fatal(err)
			return err, driver
		}
	} else {
		log.Fatal(err)
		return err, driver
	}
}
