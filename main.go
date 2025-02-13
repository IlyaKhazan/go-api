package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Flight struct {
	// переделать нейминг (camel case)
	Flight_id        int    `gorm:"flight_id" json:"id"`
	Destination_from string `json:"destination_from"`
	Destination_to   string `json:"destination_to"`
}

// передалть в БД
var flights []Flight

// var conn *pgx.Conn

// exmaple connString postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/postgres
// реализовать func getConnect(connString string) *pgx.Conn, error

// func insertDBFlight()
// func updateDBFlight()
// func getDBFlight(id)
// func getAllDBFlight(id)
// func deleteDBFlight()

func main() {
	//conn, err := getConnect
	flights = make([]Flight, 0)
	flights = append(flights, Flight{Flight_id: 1, Destination_from: "A", Destination_to: "B"})

	// вынесети в func getRounter
	router := mux.NewRouter()
	router.HandleFunc("/flights", getFlights).Methods("GET")
	router.HandleFunc("/flights/{id}", getFlight).Methods("GET")
	router.HandleFunc("/flights", createFlight).Methods("POST")
	router.HandleFunc("/flights/{id}", updateFlight).Methods("PUT")
	router.HandleFunc("/flights/{id}", deleteFlight).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getFlights(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(flights)
}

func getFlight(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO:
	// flight, err := getDBFlight()

	// for _, flight := range flights {
	// 	if flight.Flight_id == id {
	// 		json.NewEncoder(w).Encode(flight)
	// 		return
	// 	}
	// }
	http.NotFound(w, r)
}

func createFlight(w http.ResponseWriter, r *http.Request) {
	var flight Flight
	err := json.NewDecoder(r.Body).Decode(&flight)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	flight.Flight_id = len(flights) + 1
	flights = append(flights, flight)

	json.NewEncoder(w).Encode(flight)
}

func updateFlight(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var updatedFlight Flight

	err = json.NewDecoder(r.Body).Decode(&updatedFlight)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, flight := range flights {
		if flight.Flight_id == id {
			flights[i] = updatedFlight
			json.NewEncoder(w).Encode(updatedFlight)
			return
		}
	}
	http.NotFound(w, r)
}

func deleteFlight(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, flight := range flights {
		if flight.Flight_id == id {
			flights = append(flights[:i], flights[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}
