# Stage 1: Build the application
FROM golang:1.22.4 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Stage 2: Use Debian Bookworm (compatible glibc version) and install curl
FROM debian:bookworm-slim

WORKDIR /app

# Install curl for testing/debugging
RUN apt-get update && apt-get install -y curl && apt-get clean

# Copy the application binary and data files
COPY --from=build /app/main .
COPY --from=build /app/StoreMaster.csv .

EXPOSE 8080

# Run the application
CMD ["./main"]

