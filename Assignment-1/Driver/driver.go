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
	router.HandleFunc("/drivers/{email}", GetDrivers).Methods("GET")
	router.HandleFunc("/drivers/available/{available}", GetDriver).Methods("GET")
	router.HandleFunc("/drivers/{email}/tripsDriver", DriverGetTrip).Methods("GET")
	router.HandleFunc("/drivers", CreateDriver).Methods("POST")
	router.HandleFunc("/driver/{email}", UpdateDriver).Methods("PUT")
	router.HandleFunc("/driver/available/{email}", UpdateDriverAvilableStatusToFalse).Methods("PUT")
	router.HandleFunc("/drivers/starttrip/{email}", AcceptTrips).Methods("PUT")
	router.HandleFunc("/drivers/endtrip/{email}", EndTrips).Methods("PUT")
	router.HandleFunc("/drivers/{email}", DeleteDriver).Methods("Delete")

	log.Fatal(http.ListenAndServe(":50002", router))

}

const UpdateTripsStatusToProcessing = "http://localhost:50003//trips/driver"

//Initiate start trips and end trips and update the TripsStatus in trips.go3
func main() {

	initialMigration()
	initializeRouter()

}

//Driver Struct
type Driver struct {
	DriverID    int    `gorm:"primaryKey"`
	FirstName   string `json:"firstname" validate:"required"`
	Lastname    string `json:"lastname" validate:"required"`
	Phonenumber int    `json:"phonenumber" validate:"required" `
	Email       string `json:"email" validate:"required,email"`
	IcNum       string `json:"icnum" validate:"required"gorm:"<-:create"`
	LicenseNum  string `json:"licensenum" validate:"required" `
	Available   bool   `gorm:"type:bool; json:"available" validate:"required"`
}

//Trip Struct
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

var DB *gorm.DB
var err error

//Database Connections string
const ADB = "root:00Nordic00@tcp(127.0.0.2:3306)/assignment1?charset=utf8mb4&parseTime=True&loc=Local"

func DriverGetTrip(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	email := params["email"]
	if err != nil {
		//if user is not found
		fmt.Fprintf(w, "{The email you entered does not registered}")
	} else {
		err, trips := GetDriverTrip(email)
		if err != nil {
			fmt.Print("no trips")
			return
		} else {

			//gettripsbyemail := DB.Find(&Passenger{}).Where("email = ?", trip.TripPassengeremail).Order("startDate")
			json.NewEncoder(w).Encode(trips)
		}
	}
}

//Api call for get trips from trips.go
func GetDriverTrip(trip_driveremail string) (error, Trip) {
	url := UpdateTripsStatusToProcessing + "/" + trip_driveremail
	var trips Trip
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

//This fucntions for driver end-trips and call the trips update fucntions to  completed. and change the driver avilabilie status to true
func EndTrips(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var driver Driver

	//email := params["email"]
	err := DB.Where("email = ?", params["email"]).First(&driver).Error
	if err != nil {
		//if user is not found
		fmt.Fprintf(w, "The driver you entered is not registered")
	} else {

		params := mux.Vars(router)
		trip_driveremail := params["email"]
		err := UpdateTripStatusToCompleted(trip_driveremail)
		if err != nil {

			fmt.Fprintf(w, "Driver dont have any trips")
			return

		} else {
			driveremail := driver.Email
			fmt.Println(driveremail)
			DB.Model(&driver).Update("available", true)
			DB.Save(&driver)
			fmt.Fprintf(w, "Successfully update trips status")
			return
		}
	}
}

//This is url for from trip.go to update the tripstatus to complted
const UpdateTripsToCompleted = "http://localhost:50003/trips/completed"

//This functions is api call from trip.go to change the trips status to complted
func UpdateTripStatusToCompleted(trip_driveremail string) error {
	request, err := http.NewRequest(http.MethodPut,
		UpdateTripsToCompleted+"/"+trip_driveremail, nil)

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

//This fucntions for driver accept trips and call the trips updatestatustoprocess fucntions to update trips status to on process
func AcceptTrips(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var driver Driver

	//email := params["email"]
	err := DB.Where("email = ?", params["email"]).First(&driver).Error
	if err != nil {
		//if user is not found
		fmt.Fprintf(w, "The email you entered is not registered")
	} else {
		params := mux.Vars(router)
		trip_driveremail := params["email"]
		err := UpdateTripStatusToProcessing(trip_driveremail)
		if err != nil {

			fmt.Fprintf(w, "Driver dont have any trips")
			return

		} else {
			//store the email
			driveremail := driver.Email
			fmt.Println(driveremail)
			fmt.Fprintf(w, "Successfully update trips status")
			return
		}
	}

}

const UpdateTripsToProcess = "http://localhost:50003/trips/processing"

//This functions is api call from trip.go to chang the trips status to processing
func UpdateTripStatusToProcessing(trip_driveremail string) error {
	request, err := http.NewRequest(http.MethodPut,
		UpdateTripsToProcess+"/"+trip_driveremail, nil)

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

//Initial migrationg database and automigrate
func initialMigration() {
	DB, err = gorm.Open(mysql.Open(ADB), &gorm.Config{})
	/*Here is to to check if the connections to the database*/

	if err != nil {
		fmt.Println(err.Error())
		panic("cant conenct to the Database Please check the coneections strings")
	}
	DB.AutoMigrate(&Driver{})
}

//Get driver informations
func GetDrivers(w http.ResponseWriter, router *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var driver Driver
	err := DB.Where("email = ?", params["email"]).First(&driver).Error
	if err != nil {
		//if user is not found
		fmt.Fprintf(w, "The email you enter is not registered")
		return
	} else {
		json.NewEncoder(w).Encode(driver)
	}

}

//get driver by availablity status to true
func GetDriver(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//arams := mux.Vars(router)
	var driver Driver
	err := DB.Where("available = ?", true).Find(&driver).Error
	//err := DB.Where(&Driver{}, "available=?", true).First(&driver).Error
	if err != nil {
		//if user is not found
		fmt.Fprintf(w, "{}")
	} else {
		fmt.Print("available: ", driver.DriverID)
		json.NewEncoder(w).Encode(driver)
	}
}

//create drivers fucntions
func CreateDriver(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var driver Driver
	var dbdriver Driver
	json.NewDecoder(router.Body).Decode(&driver)
	//to validate the inpute must be string not empty
	validate := validator.New()
	err2 := validate.Struct(driver)
	if err2 != nil {
		fmt.Fprint(w, err2.Error())
		return
	}
	//Validate duplications of email address
	err := DB.Where("email = ?", driver.Email).First(&dbdriver).Error
	if err == nil {
		fmt.Fprint(w, "  The email you enter is already registerd")
		return
	}
	//Validate duplications of Phone number
	err3 := DB.Where("phonenumber = ?", driver.Phonenumber).First(&dbdriver).Error
	if err3 == nil {
		fmt.Fprintf(w, "  The Phone number you enter is  already registerd    ")
		return
	}

	err4 := DB.Where("icnum = ?", driver.IcNum).First(&dbdriver).Error
	if err4 == nil {
		fmt.Fprintf(w, "  The IC number you enter is  already registerd    ")
		return
	}

	err5 := DB.Where("licensenum = ?", driver.LicenseNum).First(&dbdriver).Error
	if err5 == nil {
		fmt.Fprintf(w, "  The License number you enter is  already registerd    ")
		return
	}
	DB.Create(&driver)
	json.NewEncoder(w).Encode(driver)
}

//update driver info
func UpdateDriver(w http.ResponseWriter, router *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var driver Driver

	if err != nil {
		fmt.Printf("  The email you enter is not registerd  ")
		return
	} else {
		json.NewDecoder(router.Body).Decode(&driver)
		DB.Model(&Driver{}).Where("email=?", params["email"]).Updates(driver)

		var newDriver Driver
		DB.Where("email=?", params["email"]).First(&newDriver)
		json.NewEncoder(w).Encode(newDriver)
		fmt.Printf("  Successfully update your account ")
	}
}

//This fucntions is to change the driver avilable status to false
func UpdateDriverAvilableStatusToFalse(w http.ResponseWriter, router *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var driver Driver

	err := DB.Where("email = ?", params["email"]).First(&driver).Error
	fmt.Println(params["email"])
	if err != nil {
		fmt.Printf("  The email you enter is not registerd  ")

	} else {
		json.NewDecoder(router.Body).Decode(&driver)
		fmt.Println(params)
		DB.Model(&Driver{}).Where("email = ?", params["email"]).Update("available", false)
		json.NewEncoder(w).Encode(driver)
		fmt.Printf("  Successfully update your account ")
	}
}

//Trips.go will call this api. so the driver status can change the avilable status
func UpdateDriverAvail(w http.ResponseWriter, router *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var driver Driver

	err := DB.Where("driver_id = ?", params["id"]).First(&driver).Error

	if err != nil {
		fmt.Printf(" No driver with that id  ")

	} else {
		json.NewDecoder(router.Body).Decode(&driver)
		DB.Save(&driver)
		json.NewEncoder(w).Encode(driver)
		fmt.Printf("  Successfully update your account ")

	}
}

func DeleteDriver(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var driver Driver

	err := DB.Where("email = ?", params["email"]).First(&driver).Error

	if err != nil {
		fmt.Printf("  The email you enter is not registerd  ")

	} else {
		json.NewDecoder(router.Body).Decode(&driver)
		json.NewEncoder(w).Encode("You are unable to Deleted the Account!")
	}
}
