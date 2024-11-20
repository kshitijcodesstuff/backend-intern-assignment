package models

import (
	"errors"
	"sync"
	"sync/atomic"
)

type JobRequest struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}

type Visit struct {
	StoreID   string   `json:"store_id"`
	ImageURLs []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}

var (
	jobs      = make(map[int]*Job)
	jobsMutex sync.Mutex
	jobIDSeq  int32 = 1
)

type Job struct {
	ID      int
	Request JobRequest
	Status  string
	Errors  []JobError
	Results []ImageResult
}

type JobError struct {
	StoreID string `json:"store_id"`
	Error   string `json:"error"`
}

type ImageResult struct {
	StoreID   string
	ImageURL  string
	Perimeter int
}

// CreateJob creates a new job and returns its ID
func CreateJob(req JobRequest) int {
	jobID := int(atomic.AddInt32(&jobIDSeq, 1)) - 1

	jobsMutex.Lock()
	jobs[jobID] = &Job{
		ID:      jobID,
		Request: req,
		Status:  "ongoing",
	}
	jobsMutex.Unlock()

	return jobID
}

// FetchJob retrieves a job by ID
func FetchJob(jobID int) (*Job, error) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	job, exists := jobs[jobID]
	if !exists {
		return nil, errors.New("job not found")
	}
	return job, nil
}

// AddJobError adds an error to a job
func AddJobError(jobID int, storeID, errMsg string) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	job := jobs[jobID]
	job.Errors = append(job.Errors, JobError{StoreID: storeID, Error: errMsg})
}

// FailJob sets the job status to "failed"
func FailJob(jobID int) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	job := jobs[jobID]
	job.Status = "failed"
}

// CompleteJob sets the job status to "completed"
func CompleteJob(jobID int) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	job := jobs[jobID]
	job.Status = "completed"
}

// StoreImageResult stores the result of image processing
func StoreImageResult(jobID int, storeID, imageURL string, perimeter int) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	job := jobs[jobID]
	job.Results = append(job.Results, ImageResult{
		StoreID:   storeID,
		ImageURL:  imageURL,
		Perimeter: perimeter,
	})
}

// GetJobStatus returns the status and errors of a job
func GetJobStatus(jobID int) (string, []JobError, error) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	job, exists := jobs[jobID]
	if !exists {
		return "", nil, errors.New("job not found")
	}

	return job.Status, job.Errors, nil
}
