
<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https//github.com/sahlahi/ETI-Assignment1">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>

<h3 align="center">ETI Assignment 1</h3>

  <p align="center">
    QualityRides
    <br />
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project
About QualityRides platform

Through QualityRides, the project connects millions of passengers with private hire, taxi, and coach drivers. Passengers can simply enter their current position and desired destination.


### Built With

* [golang](https://go.dev/)
* [MySql](https://www.mysql.com/)


<!-- GETTING STARTED -->
## Design consideration of your microservice

The first thing to keep in mind is that microservices are only designed to have one duty and are loosely coupled, so each service may be upgraded or restarted independently. In this case, each microservice would be responsible for (creating, reading, updating, and deleting) their specific data type from their databases.

My backend was divided into three independent microservices: passenger, driver, and trip.

This is due to the fact that each service can serve its own purpose without interfering with the others. For example, make a new passenger. It only need Passenger data.

Second, whenever a Passenger create a new trip, it requires the use of the Trips API as well as data from both the driver and the passenger. From here, I can easily construct a new trip function based on the required input (Pick up and Drop off locations). Then, I can make an api call to retrieve available driver and passenger information from the driver api and passenger api, and then add the data to new trips database. Third, when the driver and passenger want to view their trips, I can easily create functions to retrieve trips in the passenger and driver api and create an api driver and passenger api to get the trips data from trip api based on users who have trips in the trips database.

If one of the services is unavailable or experiencing a problem, the driver,passenger and trips are separated. It is simple for me to determine the source of the error and begin troubleshooting.


### Dabase


Mysql is the database I use. The reason for this is because mysql includes an optional component called the MySQL HTTP Plugin. This plugin provides direct access to MySQL through a REST over HTTP interface, removing the requirement for a middle-tier server or database-specific drivers.

### Gorm
To handle all Mysql queries, such as creating table ,column and POST GET PUT DELETE . I'm using the Gorm package so i don't need to write any querys

What is the purpose of Gorm? 

GORM provides CRUD operations and can also be used for initial migration and schema construction. GORM also excels in its extendability with native plugin support, reliance on testing, and one-to-one and one-to-many group of associations.

###  MUX
The term mux is an abbreviation for "HTTP request multiplexer." ServeMux, mux, mux, mux, mux, mux, mux, mux, mux, mux, Incoming requests are compared to a list of registered routes, and the route that matches the URL or other conditions is called.

What's the point of Gorilla mux?
The gorilla/mux package includes a request router and dispatcher for routing inbound requests to the appropriate handler. The term mux is an abbreviation for "HTTP request multiplexer."
### Validations


To validate user input for the creating new passengers, drivers, and trips. I'm use the Validator Package.

What exactly is the Validator Package?

Package validator uses tags to implement value validations for structs and individual fields. It also supports Cross-Field and Cross-Struct validation for nested structs and can dive into arrays and maps of any type.

### Database Design 

![db drawio](https://user-images.githubusercontent.com/35702954/145712619-542eb1eb-448f-4aef-b3aa-4d5f3f129fe3.png)



### Microservice Design
![db-Page-2 drawio](https://user-images.githubusercontent.com/35702954/145712620-1282cf75-fb3c-4ec5-8e6b-fb41f0eadbed.png)


Architecture diagram
### Installation

1. Clone the repo git clone https://github.com/sahlahi/ETI-Assignment1 To be Advise Install the following libraries for each microservice: Trips, Passenger, and Driver

2. Install the following packages 
    go get -u github.com/gorilla/mux
    
    go get -u gorm.io/gorm
    
    go get -u gorm.io/driver/mysql
    
    go get github.com/go-playground/validator/v10
    
    go get "github.com/go-sql-driver/mysql"
    

3. Setup Mysql Database 

   Launch the MySQL Workbench, create a new MySQL Connection, and the following window will appear. Type the following command        in the SQL File 1 window and then click the lightning icon

      CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
      GRANT ALL ON *.* TO 'user'@'localhost'
      
      ![image](https://user-images.githubusercontent.com/35702954/145713979-edc22647-72cd-482c-b8df-412bd6e6212d.png)
      
 4. Create Database enter the following command
     
     CREATE database assignment1;  USE assignment1;
     
 5. Connections string

    once you create the database schema change the connections string as shown below for each micro service 
    const ADB = "root:00Nordic00@tcp(127.0.0.2:3306)/assignment1charset=utf8mb4&parseTime=True&loc=Local"
    
    root:00Nordic00@tcp = your user account Change the user and password to the user you just create. example user:password@tcp

    (127.0.0.2:3306) here is the database IP. change it accordingly when you create. 
    
    by default the if would be like this (127.0.0.1:3306)

    /assignment1= is the database nam, change the name accordingly to your newly created database  
    
    Your connections string should look like this. Change the connections string to each microservices (Passenger,Driver,Trip)
    
    const ADB = "user:password@tcp (127.0.0.1:3306)/assignment1?charset=utf8mb4&parseTime=True&loc=Local

    Note: Do not remove ?charset=utf8mb4&parseTime=True&loc=Local




<!-- USAGE EXAMPLES -->
## How to run the program

Assignment-1 > Passenger > go mod init tidy > go run main.go

Driver

Assignment-1 > Driver > go mod init tidy > go run driver.go

Trips

Assignment-1 > Trips > go mod init tidy > go run trip.go

Frontend

Assignment-1 > frontend > go run main.go



<!-- ROADMAP -->
## Roadmap

- [] Create Databas Schema and the account
- [] Create the microservices (Passenger , Trip , driver)
- [] Create front end 

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- CONTACT -->
## Contact

Your Name - [@Email Addres ] - isyhak98@gmail.com



