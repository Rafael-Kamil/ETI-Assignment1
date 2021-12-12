package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Passenger struct {
	PassengerID int    `gorm:"primaryKey"`
	FirstName   string `json:"firstname" validate:"required"`
	Lastname    string `json:"lastname" validate:"required"`
	Phonenumber int    `json:"phonenumber" validate:"required" `
	Email       string `json:"email" validate:"required,email"`
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

var passengerEmail string
var DriverEmail string

//
func main() {
	fmt.Println("[1] Login Passenger")
	fmt.Println("[2] Login Driver")
	fmt.Println("[3] Create Passenger")
	fmt.Println("[4] Create Driver")
	fmt.Println("[0] Exit")

	var userInput string
	fmt.Scanln(&userInput)

	if userInput == "1" {
		ValidatePassenger()
	} else if userInput == "2" {
		ValidateDriver()

	} else if userInput == "3" {
		createPassenger()
	} else if userInput == "4" {
		createdriver()
	} else if userInput == "0" {
		return
	}

}

const DriverUrl = "http://localhost:50002/drivers"

//api call for passger.go to Validate Passenger email
func DriverApi(email string) (error, Driver) {
	url := DriverUrl + "/" + email
	var driver Driver
	if resp, err := http.Get(url); err == nil {
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

//Fucntions to  validate the driver account
func ValidateDriver() {
	fmt.Println("Enter your Email:")
	var credemail string
	fmt.Scanln(&credemail)
	//call the Get driver by email api to validate registered driver
	err, driver := DriverApi(credemail)
	if err != nil {
		fmt.Println("The user you enter is not registered")
		return
	} else {
		fmt.Println("Successfully loggin")
		DriverEmail = driver.Email

		fmt.Println("[1] update Info")
		fmt.Println("[2] Get Driver info")
		fmt.Println("[3] Get All trips")
		fmt.Println("[4] Start Trips")
		fmt.Println("[5] End Trips")
		fmt.Println("[6] Delete Driver")
		fmt.Println("[0] Exit")
		var userInput string
		fmt.Scanln(&userInput)

		if userInput == "1" {
			updatedriverinfo(credemail)
		} else if userInput == "2" {
			fmt.Println(driver)
		} else if userInput == "3" {
			err, trips := GetDriverTripAPi(credemail)
			if err != nil {
				fmt.Println("The user you enter is not registered")
				return
			} else {
				//for each loop, to print each trips that passenger have
				for _, trip := range trips {
					fmt.Println("DriverID: ", trip.DriverID)
					fmt.Println("PassengerID: ", trip.PassengerID)
					fmt.Println("PickipPoint: ", trip.PickUpPoint)
					fmt.Println("DropoffPoint: ", trip.DropoffPoint)
					fmt.Println("Passenger Eamil: ", trip.TripPassengeremail)
					fmt.Println("Trip Status: ", trip.TripStatus)

				}
			}

		} else if userInput == "4" {
			DriverAcceptripApi(credemail)

		} else if userInput == "5" {
			DriverEndtripApi(credemail)

		} else if userInput == "6" {
			DeleteDriverAPI(credemail)

		} else if userInput == "7" {
			return
		}

	}

}

const DriverAcceptripUrl = "http://localhost:50002/drivers/starttrip"

//api call fom driver.go to start driver trips
func DriverAcceptripApi(credemail string) error {
	request, err := http.NewRequest(http.MethodPut, DriverAcceptripUrl+"/"+credemail, nil)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
	return nil
}

const DriverEndtripUrl = "http://localhost:50002/drivers/endtrip"

//api call from driver.go to end the driver trips
func DriverEndtripApi(credemail string) error {
	request, err := http.NewRequest(http.MethodPut, DriverEndtripUrl+"/"+credemail, nil)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
	return nil
}

//this is user driver input to update driver account
func updatedriverinfo(credemail string) {
	fmt.Println("Enter your First name:")
	var firstName string
	fmt.Scanln(&firstName)

	fmt.Println("Enter your Last Name:")
	var lastName string
	fmt.Scanln(&lastName)

	fmt.Println("Enter your Phone Number:")
	var phoneNumberString string
	fmt.Scanln(&phoneNumberString)
	phoneNumber, _ := strconv.Atoi(phoneNumberString)

	fmt.Println("Enter your Email:")
	var email string
	fmt.Scanln(&email)

	fmt.Println("Enter License number :")
	var licensenum string
	fmt.Scanln(&licensenum)

	//Update driver api
	err := UpdateDriverInfo(credemail, firstName, lastName, phoneNumber, email, licensenum)
	if err != nil {
		fmt.Println("Error creating user")
	} else {
		fmt.Println("User successfully Updated")
	}
}

//this part is to stored user input to update driver account
func UpdateDriverInfo(credemail, firstName string, lastName string, phoneNumber int, email string, licensenum string) error {
	UpdateDriver := Driver{
		FirstName:   firstName,
		Lastname:    lastName,
		Phonenumber: phoneNumber,
		Email:       email,
		LicenseNum:  licensenum,
	}
	err := UpdateDriverInfoAPI(credemail, UpdateDriver)
	if err != nil {

		fmt.Println("Error updating user")
		return err
	} else {
		fmt.Println("User successfully updated")

	}
	return err
}

const UpdateDriverInfoUrl = "http://localhost:50002/driver"

//api call from passenger.go to create passenger
func UpdateDriverInfoAPI(credemail string, UpdateDriver Driver) error {
	jsonDriver, _ := json.Marshal(UpdateDriver)
	fmt.Println(credemail)
	request, err := http.NewRequest(http.MethodPut, UpdateDriverInfoUrl+"/"+credemail, bytes.NewBuffer([]byte(jsonDriver)))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
	return nil
}

//Get Drver trips from trips.go
const GetDriverTripUrl = "http://localhost:50003/trips/driver"

func GetDriverTripAPi(trip_driveremail string) (error, []Trip) {
	url := GetDriverTripUrl + "/" + trip_driveremail
	var trip []Trip
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {

			errDecode := json.Unmarshal([]byte(body), &trip)
			fmt.Println()
			if errDecode != nil {
				return err, trip
			} else {
				return err, trip
			}
		} else {
			log.Fatal(err)
			return err, trip
		}
	} else {
		log.Fatal(err)
		return err, trip
	}
}

//this is user input to create new driver
func createdriver() {
	fmt.Println("Enter your First name:")
	var firstName string
	fmt.Scanln(&firstName)

	fmt.Println("Enter your Last Name:")
	var lastName string
	fmt.Scanln(&lastName)

	fmt.Println("Enter your Phone Number:")
	var phoneNumberString string
	fmt.Scanln(&phoneNumberString)
	phoneNumber, _ := strconv.Atoi(phoneNumberString)

	fmt.Println("Enter your Email:")
	var email string
	fmt.Scanln(&email)

	fmt.Println("Enter Ic number :")
	var IcNum string
	fmt.Scanln(&IcNum)

	fmt.Println("Enter license  number :")
	var licensenum string
	fmt.Scanln(&licensenum)

	available, _ := strconv.ParseBool("true")

	//calling the stored user input fucntions
	err := createDriver(firstName, lastName, phoneNumber, email, IcNum, licensenum, available)
	if err != nil {
		fmt.Println("Error creating user")
	} else {
		fmt.Println("User successfully created")
	}
}

//this is to store createpassenger input to crea new passenger
func createDriver(firstName string, lastName string, phoneNumber int, email string, IcNum string, licensenum string, available bool) error {
	newDriver := Driver{
		FirstName:   firstName,
		Lastname:    lastName,
		Phonenumber: phoneNumber,
		Email:       email,
		IcNum:       IcNum,
		LicenseNum:  licensenum,
		Available:   available,
	}
	err := CreateDriverApi(newDriver)
	if err != nil {

		fmt.Println("Error creating user")
		return err
	} else {
		fmt.Println("User successfully created")

	}
	return err
}

//create new driver api from driver.go
const CreateDriverURL = "http://localhost:50002/drivers"

//api call from driver.go to create new driver
func CreateDriverApi(newDriver Driver) error {
	jsonDriver, _ := json.Marshal(newDriver)

	response, err := http.Post(CreateDriverURL,
		"application/json", bytes.NewBuffer([]byte(jsonDriver)))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
	return err
}

const DeleteDriverURL = "http://localhost:50002/drivers"

//api call from driver.go to delete driver
func DeleteDriverAPI(credemail string) error {
	request, err := http.NewRequest(http.MethodDelete, DeleteDriverURL+"/"+credemail, nil)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
	return nil
}

//==================End of driver Manu Sections===============================

// ==================Passenger Menu Sections========================
const PassengerURL = "http://localhost:50001/passengers"

//api call for passger.go to Validate Passenger email
func GetPassengerIdByEmail(email string) (error, Passenger) {
	url := PassengerURL + "/" + email
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

//Fucntions to  validate the passenger and passenger manu
func ValidatePassenger() {
	fmt.Println("Enter your Email:")
	var credemail string
	fmt.Scanln(&credemail)
	//
	err, passenger := GetPassengerIdByEmail(credemail)
	if err != nil {
		fmt.Println("The user you enter is not registered")
		return
	} else {
		fmt.Println("Successfully loggin")
		passengerEmail = passenger.Email

		fmt.Println("[1] update Info")
		fmt.Println("[2] Get Passenger info")
		fmt.Println("[3] Get All trips")
		fmt.Println("[4] Create new trip")
		fmt.Println("[5] Delete Accounts")
		fmt.Println("[0] Exit")
		var userInput string
		fmt.Scanln(&userInput)

		if userInput == "1" {
			updateinfo(credemail)
		} else if userInput == "2" {
			fmt.Println(passenger)
		} else if userInput == "3" {
			err, trips := GetPassengerTripsAPI(credemail)
			if err != nil {
				fmt.Println("The user you enter is not registered")
				return
			} else {
				//for each loop, to print each trips that passenger have
				for _, trip := range trips {
					fmt.Println("DriverID: ", trip.DriverID)
					fmt.Println("PassengerID: ", trip.PassengerID)
					fmt.Println("PickipPoint: ", trip.PickUpPoint)
					fmt.Println("DropoffPoint: ", trip.DropoffPoint)
					fmt.Println("Passenger Eamil: ", trip.TripPassengeremail)
					fmt.Println("Trip Status: ", trip.TripStatus)
				}
			}

		} else if userInput == "4" {
			createTrips(credemail)

		} else if userInput == "5" {
			DeletePassengerAPI(credemail)

		} else if userInput == "0" {
			return
		}

	}

}

//user input for update passenger account
func updateinfo(credemail string) {
	fmt.Println("UpdateFucntions")

	fmt.Println("Enter your First name:")
	var firstName string
	fmt.Scanln(&firstName)

	fmt.Println("Enter your Last Name:")
	var lastName string
	fmt.Scanln(&lastName)

	fmt.Println("Enter your Phone Number:")
	var phoneNumberString string
	fmt.Scanln(&phoneNumberString)
	phoneNumber, _ := strconv.Atoi(phoneNumberString)

	fmt.Println("Enter your Email:")
	var email string
	fmt.Scanln(&email)

	//CreatePassengerAPI
	err := UpdateInfo(credemail, firstName, lastName, phoneNumber, email)
	if err != nil {
		fmt.Println("Error creating user")
	} else {
		fmt.Println("User successfully Updated")
	}
}

//this part is to stored passenger input for update user account
func UpdateInfo(credemail string, firstName string, lastName string, phoneNumber int, email string) error {
	updatePassenger := Passenger{
		FirstName:   firstName,
		Lastname:    lastName,
		Phonenumber: phoneNumber,
		Email:       email,
	}
	err := UpdatePassengerInfoAPI(credemail, updatePassenger)
	if err != nil {

		fmt.Println("Error updating user")
		return err
	} else {
		fmt.Println("User successfully updated")

	}
	return err
}

const GetPassengerTrips = "http://localhost:50003/trips"

//Calling api from trip.go to get passenger trips
func GetPassengerTripsAPI(trip_passengeremail string) (error, []Trip) {
	url := GetPassengerTrips + "/" + trip_passengeremail
	var trip []Trip
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {

			errDecode := json.Unmarshal([]byte(body), &trip)
			fmt.Println()
			if errDecode != nil {
				return err, trip
			} else {
				return err, trip
			}
		} else {
			log.Fatal(err)
			return err, trip
		}
	} else {
		log.Fatal(err)
		return err, trip
	}
}

const UpdatePassengerInfoUrl = "http://localhost:50001/passengers"

//api call from passenger.go to get Reigistered Passenger info
func UpdatePassengerInfoAPI(credemail string, updatePassenger Passenger) error {
	jsonPassenger, _ := json.Marshal(updatePassenger)
	fmt.Println(credemail)
	request, err := http.NewRequest(http.MethodPut, UpdatePassengerInfoUrl+"/"+credemail, bytes.NewBuffer([]byte(jsonPassenger)))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
	return nil
}

//this is passenger input to create new passenger
func createPassenger() {
	fmt.Println("Enter your First name:")
	var firstName string
	fmt.Scanln(&firstName)

	fmt.Println("Enter your Last Name:")
	var lastName string
	fmt.Scanln(&lastName)

	fmt.Println("Enter your Phone Number:")
	var phoneNumberString string
	fmt.Scanln(&phoneNumberString)
	phoneNumber, _ := strconv.Atoi(phoneNumberString)

	fmt.Println("Enter your Email:")
	var email string
	fmt.Scanln(&email)

	err := Createpassenger(firstName, lastName, phoneNumber, email)
	if err != nil {
		fmt.Println("Error creating user")
	} else {
		fmt.Println("User successfully created")
	}
}

//this is to store createpassenger input to crea new passenger
func Createpassenger(firstName string, lastName string, phoneNumber int, email string) error {
	newPassenger := Passenger{
		FirstName:   firstName,
		Lastname:    lastName,
		Phonenumber: phoneNumber,
		Email:       email,
	}
	err := CreatePassengerAPI(newPassenger)
	if err != nil {

		fmt.Println("Error creating user")
		return err
	} else {
		fmt.Println("User successfully created")

	}
	return err
}

const CreatePassengerURL = "http://localhost:50001/passengers"

//api call from passenger.go to create passenger
func CreatePassengerAPI(newPassenger Passenger) error {
	jsonPassenger, _ := json.Marshal(newPassenger)

	response, err := http.Post(CreatePassengerURL,
		"application/json", bytes.NewBuffer([]byte(jsonPassenger)))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
	return err
}

//this is passenger input for passenger to create new trips
func createTrips(credemail string) {
	fmt.Println("Enter your pickup postal code:")
	var pickupDestinationsTostring string
	fmt.Scanln(&pickupDestinationsTostring)

	fmt.Println("Enter your Destinations postal code:")
	var dropoffdestinationsToString string
	fmt.Scanln(&dropoffdestinationsToString)

	//CreatePassengerAPI
	err := Createtrips(credemail, pickupDestinationsTostring, dropoffdestinationsToString)
	if err != nil {
		fmt.Println("Error creating user")
	} else {
		fmt.Println("User successfully created")
	}
}

//this is to store user input to create new trips
func Createtrips(credemail string, PickupDestinations string, DropOfDestinations string) error {
	newTrips := Trip{
		PickUpPoint:  PickupDestinations,
		DropoffPoint: DropOfDestinations,
	}
	err := CreateTripAPI(credemail, newTrips)
	if err != nil {

		fmt.Println("Error creating user")
		return err
	} else {
		fmt.Println("User successfully created")

	}
	return err
}

const CreateTripURL = "http://localhost:50003/trips/"

//api call from trips.go to create new trips
func CreateTripAPI(email string, newTrips Trip) error {
	jsonPassenger, _ := json.Marshal(newTrips)

	response, err := http.Post(CreateTripURL+email,
		"application/json", bytes.NewBuffer([]byte(jsonPassenger)))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()

	}
	return err
}

const DeletePassengerURL = "http://localhost:50001/passengers/"

//api call from passenger.go to Delete passenger
func DeletePassengerAPI(credemail string) error {
	request, err := http.NewRequest(http.MethodDelete, UpdatePassengerInfoUrl+"/"+credemail, nil)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
	return nil
}

//====================END of Passenger Sections========================================
