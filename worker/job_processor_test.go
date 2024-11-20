package worker

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"testing"

	"backend-intern-assignment/models"
)

// MockTransport is a custom implementation of http.RoundTripper
type MockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

// Helper to create a mock image
func createMockImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 100, 200))
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, nil)
	return buf.Bytes()
}

// Helper to initialize StoreMaster
func initTestStoreMaster() {
	models.InitTestStoreMaster()
}

// Test cases for ProcessJob
func TestProcessJob(t *testing.T) {
	// Replace the HTTP client transport with a mock
	mockTransport := &MockTransport{}
	originalTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = mockTransport
	defer func() { http.DefaultClient.Transport = originalTransport }()

	// Normal case: Valid StoreID and ImageURL
	t.Run("ValidStoreIDAndImageURL", func(t *testing.T) {
		initTestStoreMaster()

		mockTransport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(createMockImage())),
			}, nil
		}

		jobRequest := models.JobRequest{
			Count: 1,
			Visits: []models.Visit{
				{
					StoreID:   "RP00001",
					ImageURLs: []string{"https://mock-url.com/image.jpg"},
					VisitTime: "2023-10-21T15:04:05Z",
				},
			},
		}
		jobID := models.CreateJob(jobRequest)
		ProcessJob(jobID)

		job, _ := models.FetchJob(jobID)
		if job.Status != "completed" {
			t.Errorf("Expected job status 'completed', got '%s'", job.Status)
		}
		if len(job.Errors) != 0 {
			t.Errorf("Expected no errors, got %d", len(job.Errors))
		}
		if len(job.Results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(job.Results))
		}
	})

	// Edge case: Invalid StoreID
	t.Run("InvalidStoreID", func(t *testing.T) {
		initTestStoreMaster()

		jobRequest := models.JobRequest{
			Count: 1,
			Visits: []models.Visit{
				{
					StoreID:   "INVALID_STORE",
					ImageURLs: []string{"https://mock-url.com/image.jpg"},
					VisitTime: "2023-10-21T15:04:05Z",
				},
			},
		}
		jobID := models.CreateJob(jobRequest)
		ProcessJob(jobID)

		job, _ := models.FetchJob(jobID)
		if job.Status != "failed" {
			t.Errorf("Expected job status 'failed' for invalid StoreID, got '%s'", job.Status)
		}
		if len(job.Errors) != 1 {
			t.Errorf("Expected 1 error for invalid StoreID, got %d", len(job.Errors))
		}
	})

	// Edge case: Empty ImageURLs
	t.Run("EmptyImageURLs", func(t *testing.T) {
		initTestStoreMaster()

		jobRequest := models.JobRequest{
			Count: 1,
			Visits: []models.Visit{
				{
					StoreID:   "RP00001",
					ImageURLs: []string{},
					VisitTime: "2023-10-21T15:04:05Z",
				},
			},
		}
		jobID := models.CreateJob(jobRequest)
		ProcessJob(jobID)

		job, _ := models.FetchJob(jobID)
		if job.Status != "failed" {
			t.Errorf("Expected job status 'failed' for empty ImageURLs, got '%s'", job.Status)
		}
		if len(job.Errors) != 1 {
			t.Errorf("Expected 1 error for empty ImageURLs, got %d", len(job.Errors))
		}
	})

	// Edge case: Image processing failure
	t.Run("ImageProcessingFailure", func(t *testing.T) {
		initTestStoreMaster()

		mockTransport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte("not-an-image"))),
			}, nil
		}

		jobRequest := models.JobRequest{
			Count: 1,
			Visits: []models.Visit{
				{
					StoreID:   "RP00001",
					ImageURLs: []string{"https://mock-url.com"},
					VisitTime: "2023-10-21T15:04:05Z",
				},
			},
		}
		jobID := models.CreateJob(jobRequest)
		ProcessJob(jobID)

		job, _ := models.FetchJob(jobID)
		if job.Status != "failed" {
			t.Errorf("Expected job status 'failed' for image processing failure, got '%s'", job.Status)
		}
		if len(job.Errors) != 1 {
			t.Errorf("Expected 1 error for image processing failure, got %d", len(job.Errors))
		}
	})

	// Edge case: Failed HTTP request
	t.Run("FailedHTTPRequest", func(t *testing.T) {
		initTestStoreMaster()

		mockTransport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewReader([]byte("server error"))),
			}, nil
		}

		jobRequest := models.JobRequest{
			Count: 1,
			Visits: []models.Visit{
				{
					StoreID:   "RP00001",
					ImageURLs: []string{"https://mock-url.com/image.jpg"},
					VisitTime: "2023-10-21T15:04:05Z",
				},
			},
		}
		jobID := models.CreateJob(jobRequest)
		ProcessJob(jobID)

		job, _ := models.FetchJob(jobID)
		if job.Status != "failed" {
			t.Errorf("Expected job status 'failed' for failed HTTP request, got '%s'", job.Status)
		}
		if len(job.Errors) != 1 {
			t.Errorf("Expected 1 error for failed HTTP request, got %d", len(job.Errors))
		}
	})
}
