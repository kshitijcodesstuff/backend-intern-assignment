package worker

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"backend-intern-assignment/models"
	"backend-intern-assignment/utils"
)

// HTTPClient is the client used for HTTP requests. It can be overridden during tests.
var HTTPClient = &http.Client{}

var randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

// ProcessJob processes images for a job
// func ProcessJob(jobID int) {
// 	startTime := time.Now()

// 	// Fetch the job from the job store
// 	job, err := models.FetchJob(jobID)
// 	if err != nil {
// 		log.Printf("Failed to fetch job: %v", err)
// 		return
// 	}

// 	var hasErrors bool

// 	// Iterate over each visit in the job request
// 	for _, visit := range job.Request.Visits {
// 		// Validate the store ID
// 		if !models.IsValidStore(visit.StoreID) {
// 			models.AddJobError(jobID, visit.StoreID, "Invalid store ID")
// 			hasErrors = true
// 			continue
// 		}

// 		// Process each image URL for the visit
// 		for _, imageURL := range visit.ImageURLs {
// 			// Download the image using the HTTP client
// 			resp, err := HTTPClient.Get(imageURL)
// 			if err != nil || resp.StatusCode != http.StatusOK {
// 				models.AddJobError(jobID, visit.StoreID, "Failed to download image")
// 				hasErrors = true
// 				continue
// 			}

// 			// Calculate the perimeter of the image
// 			perimeter, err := utils.CalculatePerimeter(resp.Body)
// 			resp.Body.Close()
// 			if err != nil {
// 				models.AddJobError(jobID, visit.StoreID, "Failed to process image")
// 				hasErrors = true
// 				continue
// 			}

// 			// Simulate GPU processing with a random delay between 100ms to 400ms
// 			delay := time.Duration(randomGenerator.Intn(301)+100) * time.Millisecond
// 			time.Sleep(delay)

// 			// Store the image result in the job
// 			models.StoreImageResult(jobID, visit.StoreID, imageURL, perimeter)
// 		}
// 	}

// 	// Update the job status based on whether there were errors
// 	if hasErrors {
// 		models.FailJob(jobID)
// 	} else {
// 		models.CompleteJob(jobID)
// 	}

//		// Log the total processing time for the job
//		totalTime := time.Since(startTime)
//		log.Printf("Job ID %d: Total processing time %v", jobID, totalTime)
//	}
func ProcessJob(jobID int) {
	startTime := time.Now()

	job, err := models.FetchJob(jobID)
	if err != nil {
		log.Printf("Failed to fetch job: %v", err)
		return
	}

	var hasErrors bool

	for _, visit := range job.Request.Visits {
		log.Printf("Processing visit for Store ID: %s", visit.StoreID)

		// Check if StoreID is valid
		if !models.IsValidStore(visit.StoreID) {
			log.Printf("Invalid Store ID: %s", visit.StoreID)
			models.AddJobError(jobID, visit.StoreID, "Invalid Store ID")
			hasErrors = true
			continue
		}

		// Check if ImageURLs is empty
		if len(visit.ImageURLs) == 0 {
			log.Printf("Empty ImageURLs for Store ID: %s", visit.StoreID)
			models.AddJobError(jobID, visit.StoreID, "No images provided for processing")
			hasErrors = true
			continue
		}

		for _, imageURL := range visit.ImageURLs {
			log.Printf("Downloading image: %s", imageURL)

			resp, err := http.Get(imageURL)
			if err != nil || resp.StatusCode != http.StatusOK {
				log.Printf("Failed to download image: %s", imageURL)
				models.AddJobError(jobID, visit.StoreID, "Failed to download image")
				hasErrors = true
				continue
			}

			log.Printf("Processing image: %s", imageURL)
			perimeter, err := utils.CalculatePerimeter(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Printf("Failed to process image: %s", imageURL)
				models.AddJobError(jobID, visit.StoreID, "Failed to process image")
				hasErrors = true
				continue
			}

			// Simulate GPU processing delay
			delay := time.Duration(randomGenerator.Intn(301)+100) * time.Millisecond
			log.Printf("Simulating GPU processing with delay: %v", delay)
			time.Sleep(delay)

			models.StoreImageResult(jobID, visit.StoreID, imageURL, perimeter)
			log.Printf("Successfully processed image: %s with perimeter: %d", imageURL, perimeter)
		}
	}

	// Mark the job status based on whether there were errors
	if hasErrors {
		log.Printf("Job ID %d: Marking job as failed", jobID)
		models.FailJob(jobID)
	} else {
		log.Printf("Job ID %d: Marking job as completed", jobID)
		models.CompleteJob(jobID)
	}

	totalTime := time.Since(startTime)
	log.Printf("Job ID %d: Total processing time %v", jobID, totalTime)
}
