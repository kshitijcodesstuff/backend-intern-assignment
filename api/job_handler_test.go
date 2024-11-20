package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"backend-intern-assignment/models"
	"backend-intern-assignment/worker"

	"github.com/gorilla/mux"
)

// Helper to set up the router
func setupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/submit/", SubmitJob).Methods("POST")
	router.HandleFunc("/api/status", GetJobStatus).Methods("GET")
	return router
}

func TestSubmitJob(t *testing.T) {
	router := setupRouter()

	// Normal case: Valid job request
	t.Run("ValidJobRequest", func(t *testing.T) {
		jobRequest := models.JobRequest{
			Count: 1,
			Visits: []models.Visit{
				{
					StoreID:   "RP00001",
					ImageURLs: []string{"https://example.com/image.jpg"},
					VisitTime: "2023-10-21T15:04:05Z",
				},
			},
		}
		payload, _ := json.Marshal(jobRequest)
		req, _ := http.NewRequest("POST", "/api/submit/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Errorf("Expected status code 201, got %d", resp.Code)
		}

		var response map[string]int
		json.Unmarshal(resp.Body.Bytes(), &response)
		if _, exists := response["job_id"]; !exists {
			t.Errorf("Response does not contain job_id")
		}
	})

	// Edge case: Invalid JSON payload
	t.Run("InvalidJSON", func(t *testing.T) {
		invalidPayload := []byte(`{ invalid json }`)
		req, _ := http.NewRequest("POST", "/api/submit/", bytes.NewBuffer(invalidPayload))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400 for invalid JSON, got %d", resp.Code)
		}
	})

	// Edge case: Mismatched count and visits
	t.Run("MismatchedCountAndVisits", func(t *testing.T) {
		jobRequest := models.JobRequest{
			Count: 2,
			Visits: []models.Visit{
				{
					StoreID:   "RP00001",
					ImageURLs: []string{"https://example.com/image.jpg"},
					VisitTime: "2023-10-21T15:04:05Z",
				},
			},
		}
		payload, _ := json.Marshal(jobRequest)
		req, _ := http.NewRequest("POST", "/api/submit/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400 for mismatched count, got %d", resp.Code)
		}
	})
}

func TestGetJobStatus(t *testing.T) {
	router := setupRouter()

	// Normal case: Valid job ID
	t.Run("ValidJobID", func(t *testing.T) {
		// Create a dummy job
		jobRequest := models.JobRequest{
			Count: 1,
			Visits: []models.Visit{
				{
					StoreID:   "RP00001",
					ImageURLs: []string{"https://example.com/image.jpg"},
					VisitTime: "2023-10-21T15:04:05Z",
				},
			},
		}
		jobID := models.CreateJob(jobRequest)
		go worker.ProcessJob(jobID)

		req, _ := http.NewRequest("GET", "/api/status?jobid="+strconv.Itoa(jobID), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", resp.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &response)
		if response["status"] != "completed" && response["status"] != "ongoing" {
			t.Errorf("Unexpected job status: %s", response["status"])
		}
	})

	// Edge case: Missing job ID parameter
	t.Run("MissingJobID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/status", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400 for missing job ID, got %d", resp.Code)
		}
	})

	// Edge case: Invalid job ID format
	t.Run("InvalidJobIDFormat", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/status?jobid=invalid", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400 for invalid job ID format, got %d", resp.Code)
		}
	})

	// Edge case: Non-existent job ID
	t.Run("NonExistentJobID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/status?jobid=99999", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400 for non-existent job ID, got %d", resp.Code)
		}
	})
}
