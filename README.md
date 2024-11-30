
# Receipt Processor

This is a REST API service that processes receipts and calculates points based on specific rules. The service provides endpoints to submit receipts and retrieve their calculated points.

## Features

- Process receipts and generate unique IDs
- Calculate points based on various receipt attributes
- In-memory storage for receipt data
- RESTful API endpoints
- Docker support

## Prerequisites

- Go 1.21 or higher (for direct execution)
- Docker (for containerized execution)

## Running the Application

### Using Docker

1. Clone the repository:
```bash
git clone https://github.com/pillows/receipt-processor
cd receipt-processor
```

2. Build the Docker image:
```bash
docker build -t receipt-processor .
```

3. Run the container:
```bash
docker run -p 8000:8000 receipt-processor
```

The service will be available at `http://localhost:8000`

### Using Go (Direct Execution)

1. Clone the repository:
```bash
git clone https://github.com/pillows/receipt-processor
cd receipt-processor
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

The service will be available at `http://localhost:8000`

## Development

### Project Structure
```
receipt-processor/
├── Dockerfile
├── docker-compose.yml
├── README.md
├── go.mod
├── go.sum
├── main.go
├── run.sh
└── models/
    └── receipt.go
└── utils/
    └── points.go
```
## Notes

- There are detailed error message descriptions but are commented out in favor of matching the api spec.