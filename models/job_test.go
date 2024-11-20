package models

import (
	"testing"
)

func TestCreateAndFetchJob(t *testing.T) {
	// Normal case: Create and fetch a valid job
	t.Run("ValidJobCreationAndFetch", func(t *testing.T) {
		jobRequest := JobRequest{
			Count: 1,
			Visits: []Visit{
				{
					StoreID:   "RP00001",
					ImageURLs: []string{"https://www.example.com/image.jpg"},
					VisitTime: "2023-10-21T15:04:05Z",
				},
			},
		}
		jobID := CreateJob(jobRequest)
		job, err := FetchJob(jobID)
		if err != nil {
			t.Errorf("Expected to fetch job without error, got %v", err)
		}
		if job.ID != jobID {
			t.Errorf("Expected job ID %d, got %d", jobID, job.ID)
		}
	})

	// Edge case: Fetch a non-existent job
	t.Run("FetchNonExistentJob", func(t *testing.T) {
		_, err := FetchJob(999) // Assuming 999 is an invalid job ID
		if err == nil {
			t.Error("Expected error for non-existent job, got none")
		}
	})

	// Edge case: Create a job with zero visits
	t.Run("JobWithZeroVisits", func(t *testing.T) {
		jobRequest := JobRequest{
			Count:  0,
			Visits: []Visit{},
		}
		jobID := CreateJob(jobRequest)
		job, err := FetchJob(jobID)
		if err != nil {
			t.Errorf("Expected to fetch job without error, got %v", err)
		}
		if job.Request.Count != 0 || len(job.Request.Visits) != 0 {
			t.Errorf("Expected job with zero visits, got count %d and %d visits", job.Request.Count, len(job.Request.Visits))
		}
	})
}

func TestAddJobError(t *testing.T) {
	// Normal case: Add a single error to a job
	t.Run("AddSingleError", func(t *testing.T) {
		jobRequest := JobRequest{
			Count:  1,
			Visits: []Visit{},
		}
		jobID := CreateJob(jobRequest)
		AddJobError(jobID, "RP00001", "Test error")
		job, _ := FetchJob(jobID)
		if len(job.Errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(job.Errors))
		}
		if job.Errors[0].Error != "Test error" {
			t.Errorf("Expected error message 'Test error', got '%s'", job.Errors[0].Error)
		}
	})

	// Edge case: Add multiple errors to a job
	t.Run("AddMultipleErrors", func(t *testing.T) {
		jobRequest := JobRequest{
			Count:  1,
			Visits: []Visit{},
		}
		jobID := CreateJob(jobRequest)
		AddJobError(jobID, "RP00001", "First error")
		AddJobError(jobID, "RP00002", "Second error")
		job, _ := FetchJob(jobID)
		if len(job.Errors) != 2 {
			t.Errorf("Expected 2 errors, got %d", len(job.Errors))
		}
		if job.Errors[1].Error != "Second error" {
			t.Errorf("Expected error message 'Second error', got '%s'", job.Errors[1].Error)
		}
	})

	// Edge case: Add an error to a non-existent job
	t.Run("AddErrorToNonExistentJob", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when adding error to non-existent job, got none")
			}
		}()
		AddJobError(999, "RP00001", "Non-existent job error")
	})
}

func TestJobStatusFunctions(t *testing.T) {
	// Normal case: Fail and then complete a job
	t.Run("FailAndCompleteJob", func(t *testing.T) {
		jobRequest := JobRequest{
			Count:  1,
			Visits: []Visit{},
		}
		jobID := CreateJob(jobRequest)
		FailJob(jobID)
		job, _ := FetchJob(jobID)
		if job.Status != "failed" {
			t.Errorf("Expected job status 'failed', got '%s'", job.Status)
		}

		CompleteJob(jobID)
		job, _ = FetchJob(jobID)
		if job.Status != "completed" {
			t.Errorf("Expected job status 'completed', got '%s'", job.Status)
		}
	})

	// Edge case: Fail a job multiple times
	t.Run("FailJobMultipleTimes", func(t *testing.T) {
		jobRequest := JobRequest{
			Count:  1,
			Visits: []Visit{},
		}
		jobID := CreateJob(jobRequest)
		FailJob(jobID)
		FailJob(jobID) // Should remain in "failed" state
		job, _ := FetchJob(jobID)
		if job.Status != "failed" {
			t.Errorf("Expected job status 'failed', got '%s'", job.Status)
		}
	})

	// Edge case: Complete a job without failing
	t.Run("CompleteJobDirectly", func(t *testing.T) {
		jobRequest := JobRequest{
			Count:  1,
			Visits: []Visit{},
		}
		jobID := CreateJob(jobRequest)
		CompleteJob(jobID)
		job, _ := FetchJob(jobID)
		if job.Status != "completed" {
			t.Errorf("Expected job status 'completed', got '%s'", job.Status)
		}
	})

	// Edge case: Change status of non-existent job
	t.Run("ChangeStatusOfNonExistentJob", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when changing status of non-existent job, got none")
			}
		}()
		FailJob(999) // Assuming 999 is an invalid job ID
	})
}
