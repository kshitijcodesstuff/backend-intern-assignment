package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-intern-assignment/api"
	"backend-intern-assignment/db"
	"backend-intern-assignment/models"

	"github.com/gorilla/mux"
)

func setupTestServer() *mux.Router {
	db.InitDB()
	models.LoadStoreMaster("StoreMaster.csv")
	r := mux.NewRouter()
	r.HandleFunc("/api/submit/", api.SubmitJob).Methods("POST")
	r.HandleFunc("/api/status", api.GetJobStatus).Methods("GET")
	return r
}

func TestServerStartup(t *testing.T) {
	router := setupTestServer()

	req, _ := http.NewRequest("GET", "/api/status?jobid=1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", resp.Code)
	}
}
