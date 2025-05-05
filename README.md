# Anomaly Detection System

A Go-based API for detecting and managing anomalies in time-series data.

## Features

- CRUD operations for anomalies
- Configurable detection rules
- RESTful API endpoints
- Super basic frontend

## Prerequisites

- Go 1.16 or higher
- Gin web framework
- Postgresql v14
- Node.js v20 or higher
- pnpm v9 or higher

## Installation

1. Clone the repository:
```bash
git clone https://github.com/ainesh01/anomaly_detection.git
cd anomaly_detection
```

2. Install dependencies:
```bash
go mod tidy && go mod vendor && go mod tidy
```

## Running the Application

To start the server:

```bash
go run cmd/api/main.go
```

To run the frontend:
```bash
cd frontend
pnpm install
pnpm dev
```

## Future Improvements

- Implement proper database integration
- Add authentication and authorization
- Implement indepth anomaly detection rules
- Add unit tests and integration tests
- Add API documentation with Swagger
- Add configuration management
- Add Docker containerization for easy deployment
- Add CI/CD pipeline

## Design Decisions

A major design decision was to use a Postgresql database to store the data. This was chosen because it is a mature and well-supported database that is easy to set up and use. It also allows for the use of the full range of SQL features, which were used to implement the more complex queries required for the anomaly detection. 

I also chose to add a standard deviation anomaly check to check for outliers in the data. This is a simple and effective way to detect anomalies in the data, and is a good starting point for the anomaly detection. It is also computationally efficient, and does not require a machine learning model. A SD greater than an absolute value of 3, should be investigated for further analysis. 

This is a very simple implementation of the anomaly detection system, and there are many ways to improve it. The code should be refactored to be more modular and easier to maintain. The anomaly detection should be automatically run when a new rule is added or new data is inserted, and only for the rule being added/changed or for the data set being inserted to reduce computational load. The code should be refactored to be more efficient, and possibly use a machine learning model for the anomaly detection, which is out of the scope of this takehome. 

AnomalyRules can only be used on three fields: `max_salary`, `min_salary`, and `company_rating` at the moment. This should be refactored to dynamically add the fields to the rule from the dataset being passed in.

I personally have little experience with frontend, and made a very simple front end using heavy AI assistance. 

I also stubbed an `AdvancedAnomalyRule` model, but did not have time to implement it. The code should be refactored to use the `AdvancedAnomalyRule` model, and the `AnomalyRule` model should be removed. This will allow for more complex anomaly detection rules to be added, and the code will be more maintainable, as well as adding severity levels to the anomaly detection. 

## New Data
New data can be POSTed to the server using the `POST /api/job-data` endpoint.

## Anomaly Rules
Anomaly rules can be POSTed to the server using the `POST /api/anomaly-rules` endpoint or via the frontend.

## Accessing the frontend
The frontend can be accessed at `http://localhost:3000/`.

## Accessing the API
The API can be accessed at `http://localhost:8080/api/`.

