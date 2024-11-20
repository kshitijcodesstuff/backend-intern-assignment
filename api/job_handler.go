package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend-intern-assignment/models"
	"backend-intern-assignment/worker"
)

// SubmitJob handles job submission
func SubmitJob(w http.ResponseWriter, r *http.Request) {
	var jobRequest models.JobRequest
	err := json.NewDecoder(r.Body).Decode(&jobRequest)
	if err != nil || jobRequest.Count != len(jobRequest.Visits) {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	jobID := models.CreateJob(jobRequest)
	go worker.ProcessJob(jobID)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"job_id": jobID})
}

// GetJobStatus retrieves the status of a job
func GetJobStatus(w http.ResponseWriter, r *http.Request) {
	jobIDParam := r.URL.Query().Get("jobid")
	jobID, err := strconv.Atoi(jobIDParam)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	status, errors, err := models.GetJobStatus(jobID)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"status": status,
		"job_id": jobID,
	}
	if status == "failed" {
		response["error"] = errors
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
