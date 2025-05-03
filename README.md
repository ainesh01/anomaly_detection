# Anomaly Detection System

A Go-based API for detecting and managing anomalies in time-series data.

## Features

- CRUD operations for anomalies
- Configurable detection rules
- RESTful API endpoints
- Health check endpoint

## Prerequisites

- Go 1.16 or higher
- Gin web framework

## Installation

1. Clone the repository:
```bash
git clone https://github.com/ainesh01/anomaly_detection.git
cd anomaly_detection
```

2. Install dependencies:
```bash
go mod tidy
```

## Running the Application

To start the server:

```bash
go run cmd/api/main.go
```

The server will start on port 8080.

## API Endpoints

### Health Check
- `GET /health` - Check server status

### Anomalies
- `GET /api/anomalies` - Get all anomalies
- `GET /api/anomalies/:id` - Get a specific anomaly


## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── models/
│   │   └── anomaly.go
│   └── handlers/
│       └── anomaly_handler.go
├── pkg/
│   └── database/
├── go.mod
└── README.md
```

## Future Improvements

- Implement database integration
- Add authentication and authorization
- Implement anomaly detection rules
- Add data ingestion endpoints
- Add frontend visualization
- Add unit tests and integration tests
- Add API documentation with Swagger
- Implement logging and monitoring
- Add configuration management
- Add Docker support 