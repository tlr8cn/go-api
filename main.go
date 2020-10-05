package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

const (
	TIMEOUT = 20
)

var (
	db *gorm.DB

	previous *Previous
	current  *Current
	next     *Next
)

type EndpointResponse struct {
	Value int `json:"value"`
}

type Previous EndpointResponse
type Current EndpointResponse
type Next EndpointResponse

func main() {
	var err error

	// Assumes existance of local MySQL database called "fibonacci"
	db, err = gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/fibonacci?charset=utf8&parseTime=True")
	fatalIf(err)

	//Based on Struct definitions, gorm will handle the schema for us
	db.AutoMigrate(&Previous{})
	db.AutoMigrate(&Current{})
	db.AutoMigrate(&Next{})

	//Create a new thread to step through the Fibonacci Sequence
	go stepThroughFibonacci()

	//Handle Requests on the main thread
	handleRequests()
}

func getPrevious(w http.ResponseWriter, r *http.Request) {
	var tmpPrevious *Previous

	if previous == nil {
		tmpPrevious = &Previous{Value: 0}
	} else {
		tmpPrevious = &Previous{}
		db.Find(tmpPrevious)
	}

	json.NewEncoder(w).Encode(tmpPrevious)
}

func getCurrent(w http.ResponseWriter, r *http.Request) {
	var tmpCurrent *Current

	if current == nil {
		tmpCurrent = &Current{Value: 0}
	} else {
		tmpCurrent = &Current{}
		db.Find(tmpCurrent)
	}

	json.NewEncoder(w).Encode(tmpCurrent)
}

func getNext(w http.ResponseWriter, r *http.Request) {
	var tmpNext *Next

	if next == nil {
		tmpNext = &Next{Value: 1}
	} else {
		tmpNext = &Next{}
		db.Find(tmpNext)
	}

	json.NewEncoder(w).Encode(tmpNext)
}

func handleRequests() {
	log.Println("Starting up...")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/previous", getPrevious)
	router.HandleFunc("/current", getCurrent)
	router.HandleFunc("/next", getNext)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.RecoveryHandler()(router)))
}

func stepThroughFibonacci() {
	startSequence()

	for {
		time.Sleep(TIMEOUT * time.Second)

		updatePrevious()
		updateCurrent()
		updateNext()

		log.Println("Values Updated:")
		log.Printf("Previous: %d\n", previous.Value)
		log.Printf("Current: %d\n", current.Value)
		log.Printf("Next: %d\n", next.Value)
	}
}

func startSequence() {
	if current != nil {
		db.Delete(current)
	}

	current = &Current{Value: 0}
	db.Create(current)

	if next != nil {
		db.Delete(next)
	}

	next = &Next{Value: 1}
	db.Create(next)

	if previous != nil {
		db.Delete(previous)
	}

	previous = &Previous{Value: 0}
	db.Create(previous)

}

func updatePrevious() {
	if previous != nil {
		db.Delete(previous)
	}
	previous = &Previous{Value: current.Value}
	db.Create(previous)
}

func updateCurrent() {
	db.Delete(current)
	current = &Current{Value: next.Value}
	db.Create(current)
}

func updateNext() {
	db.Delete(next)
	nextVal := current.Value + previous.Value
	if nextVal < 0 { //Protect against integer overflow
		startSequence()
		nextVal = 1
	}
	next = &Next{Value: nextVal}
	db.Create(next)
}

func fatalIf(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}
