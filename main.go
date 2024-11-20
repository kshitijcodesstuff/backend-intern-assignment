package main

import (
	"log"
	"net/http"

	"backend-intern-assignment/api"
	"backend-intern-assignment/db"
	"backend-intern-assignment/models"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database and preload StoreMaster data
	db.InitDB()
	models.LoadStoreMaster("StoreMaster.csv")

	// Set up router and endpoints
	r := mux.NewRouter()
	r.HandleFunc("/api/submit/", api.SubmitJob).Methods("POST")
	r.HandleFunc("/api/status", api.GetJobStatus).Methods("GET")

	// Start the server
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
