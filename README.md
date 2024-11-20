
# Kirana Backend Assignment

## **GitHub Repository**
[GitHub Repository Link](https://github.com/kshitijcodesstuff/backend-intern-assignment)

---

## **Overview**
This project is a backend system developed to handle job submissions, processes store visit data, calculates image perimeters, validates store IDs from a master list, and provides APIs for job management. It is written in **Go**, uses in-memory storage for simplicity, and is containerized using **Docker**.

---

## **Project Structure**

The project is organized into a well-structured directory hierarchy to separate concerns, making it easy to understand, maintain, and extend.

```plaintext
backend-intern-assignment/
├── main.go                      # Entry point of the application.
├── StoreMaster.csv              # CSV file containing store data with Store IDs and names.
├── api/                         # API layer for managing HTTP endpoints.
│   ├── job_handler.go           # Handles API requests for job submission and status retrieval.
│   ├── job_handler_test.go      # Unit tests for the job handler functions.
├── worker/                      # Worker layer for background job processing.
│   ├── job_processor.go         # Processes jobs asynchronously (image downloading, validation, etc.).
│   ├── job_processor_test.go    # Unit tests for the job processing logic.
├── utils/                       # Utility layer for reusable functions.
│   ├── utils.go                 # Provides functions like image perimeter calculation.
│   ├── utils_test.go            # Unit tests for utility functions.
├── models/                      # Models layer for managing data structures and logic.
│   ├── job.go                   # Models and logic for job management, including status updates.
│   ├── job_test.go              # Unit tests for job-related logic.
│   ├── store_master.go          # Logic for loading and validating store data from StoreMaster.csv.
│   ├── store_master_test.go     # Unit tests for store master functionality.
├── db/                          # Database-related setup and configuration.
│   ├── database.go              # Initializes in-memory storage for jobs and other data.
├── go.mod                       # Go module file listing dependencies for the project.
├── go.sum                       # Checksums for verifying module integrity.
├── Dockerfile                   # Dockerfile for containerizing the application.
```


---

## **Core Features**

### **1. Job Submission**
- **Endpoint**: `/api/submit/` (POST)
- **Description**: Submits a job to process store visit data and images.
- **Request Format**:
    ```json
    {
        "count": 1,
        "visits": [
            {
                "store_id": "RP00001",
                "image_url": ["https://example.com/image.jpg"],
                "visit_time": "2023-10-21T15:04:05Z"
            }
        ]
    }
    ```
- **Response**:
    - On Success:
      ```json
      {
          "job_id": 1
      }
      ```
    - On Failure (e.g., invalid input):
      ```json
      {
          "error": "Invalid request"
      }
      ```

---

### **2. Job Status Retrieval**
- **Endpoint**: `/api/status` (GET)
- **Description**: Retrieves the current status of a job by ID.
- **Query Parameter**:
    - `jobid`: The ID of the job.
- **Response**:
    - **Job Completed**:
      ```json
      {
          "status": "completed",
          "job_id": 1
      }
      ```
    - **Job Ongoing**:
      ```json
      {
          "status": "ongoing",
          "job_id": 1
      }
      ```
    - **Job Failed**:
      ```json
      {
          "status": "failed",
          "job_id": 1,
          "error": [
              {"store_id": "RP00001", "error": "Invalid Store ID"}
          ]
      }
      ```
    - **Invalid Job ID**:
      ```json
      {
          "error": "Job not found"
      }
      ```

---

### **3. Store Validation**
- **Source**: `StoreMaster.csv`
- **Description**:
  - Preloads valid store IDs from the CSV file at startup.
  - Validates store IDs during job processing.

---

### **4. Image Processing**
- **Description**:
  - Decodes image files to calculate their perimeters using dimensions.
  - Simulates GPU processing delays for realism.

---

### **5. Job Processing**
- **Description**:
  - Validates store IDs against the master list.
  - Downloads and processes images for perimeter calculation.
  - Handles errors (e.g., invalid store IDs, image download failures, or empty image lists).

---

## **Error Handling**
- **Scenarios**:
  - Invalid request payloads: Responds with `400 Bad Request`.
  - Non-existent job IDs: Returns appropriate error messages.
  - Image processing errors: Marks jobs as "failed" with detailed error descriptions.

---

## **Setup Instructions**

### **Prerequisites**
- **Go**: Version 1.22.4
- **Docker**: Installed and configured.

---

### **Run Locally**
1. Clone the repository:
    ```bash
    git clone https://github.com/your-repo/backend-intern-assignment.git
    cd backend-intern-assignment
    ```
2. Install dependencies:
    ```bash
    go mod download
    ```
3. Run the application:
    ```bash
    go run main.go
    ```
4. Access the application:
    - Submit jobs: `http://localhost:8080/api/submit/`
    - Get job status: `http://localhost:8080/api/status?jobid=1`

---

### **Run with Docker**
1. Build the Docker image:
    ```bash
    docker build -t backend-intern-assignment .
    ```
2. Run the container:
    ```bash
    docker run -p 8080:8080 backend-intern-assignment
    ```
3. Access the APIs at `http://localhost:8080`.

---

## **Testing**

### **Running Tests**
1. Run all tests:
    ```bash
    go test ./... -v
    ```
2. Run tests for a specific component:
    ```bash
    go test ./api -v
    ```

---

### **Test Coverage**

#### **API Tests**
- Valid job submission and retrieval.
- Edge cases:
  - Mismatched counts.
  - Invalid JSON payloads.
  - Non-existent or missing job IDs.

#### **Job Models Tests**
- Job creation, status updates, and error handling.
- Edge cases:
  - Jobs with no visits.
  - Updating statuses of non-existent jobs.

#### **Worker Tests**
- Simulates job processing:
  - Valid and invalid store IDs.
  - Empty or corrupted image data.
  - Network request failures.

#### **Utility Tests**
- Validates image perimeter calculations for valid and invalid images.

#### **Store Validation Tests**
- Ensures `StoreMaster.csv` is loaded correctly and validates store IDs.

---

## **Example cURL Commands**

### Submit a Job
```bash
curl -X POST http://localhost:8080/api/submit/ \
-H "Content-Type: application/json" \
-d '{
    "count": 1,
    "visits": [
        {
            "store_id": "RP00001",
            "image_url": ["https://example.com/image.jpg"],
            "visit_time": "2023-10-21T15:04:05Z"
        }
    ]
}'
```

### Retrieve Job Status
```bash
curl -X GET "http://localhost:8080/api/status?jobid=1"
```

---

## **Why This Application Stands Out**
1. **Rigorous Testing**: Extensively tested across normal and edge cases, ensuring reliability and robustness.
2. **Modular Design**: Clean separation of components for easier development and maintenance.
3. **Comprehensive Error Handling**: Provides detailed error messages for users and developers.
4. **Realistic Simulations**: Incorporates GPU delays and handles network failures gracefully.
5. **Production-Ready**: Supports asynchronous processing, error recovery, and scalability.

---

## **Future Enhancements**
1. Replace in-memory storage with a database like PostgreSQL or Redis.
2. Integrate a distributed task queue (e.g., RabbitMQ) for job management.
3. Add monitoring tools like Prometheus or Grafana for real-time insights.

---

## **Conclusion**
This backend application for **Kirana** demonstrates robust design, effective error handling, and thorough testing. It is a production-ready solution built to handle real-world scenarios and is flexible enough for future growth.
