// Package services provides core business logic for anomaly detection in job listings.
// It includes functionality for detecting both absolute and relative anomalies in job data,
// managing anomaly statistics, and persisting anomalies to the database.
package services

import (
	"fmt"
	"math"
	"time"

	"github.com/ainesh01/anomaly_detection/internal/models"
	"github.com/lib/pq"
)

// AnomalyServiceInterface defines the interface for anomaly detection and retrieval operations
type AnomalyServiceInterface interface {
	DetectAnomalies(job *models.JobData) ([]models.Anomaly, error)
	GetAnomaliesByJobID(jobID string) ([]models.Anomaly, error)
	GetAllAnomalies() ([]models.Anomaly, error)
	DetectAnomaliesForAllJobs() error
}

// AnomalyType represents the specific type of anomaly detected
type AnomalyType string

// AnomalyCategory indicates whether an anomaly is based on absolute rules or relative statistics
type AnomalyCategory string

// ComparisonOperator represents the type of comparison used in absolute anomaly detection
type ComparisonOperator string

const (
	// Anomaly types
	AnomalyTypeNullValues AnomalyType = "null_values"        // When required fields are null
	AnomalyTypeDeviation  AnomalyType = "standard_deviation" // When value deviates significantly from mean
	AnomalyTypeThreshold  AnomalyType = "threshold"          // When value exceeds a threshold
	AnomalyTypeMinSalary  AnomalyType = "min_salary"         // When value deviates significantly from mean

	// Anomaly categories
	AnomalyCategoryAbsolute AnomalyCategory = "absolute" // Based on fixed rules
	AnomalyCategoryRelative AnomalyCategory = "relative" // Based on statistical analysis

	// Comparison operators for absolute anomalies
	GreaterThan        ComparisonOperator = ">"
	GreaterThanOrEqual ComparisonOperator = ">="
	LessThan           ComparisonOperator = "<"
	LessThanOrEqual    ComparisonOperator = "<="
	Equal              ComparisonOperator = "="

	// Standard deviation threshold for anomaly detection
	StdDevThreshold = 3.0
)

// ValidOperators is a list of all valid comparison operators
var ValidOperators = []ComparisonOperator{
	GreaterThan,
	GreaterThanOrEqual,
	LessThan,
	LessThanOrEqual,
	Equal,
}

// IsValidOperator checks if the given operator is valid
func IsValidOperator(op ComparisonOperator) bool {
	for _, validOp := range ValidOperators {
		if op == validOp {
			return true
		}
	}
	return false
}

// Anomaly represents a detected anomaly in job data
type Anomaly struct {
	ID          int64              `json:"id"`
	JobID       string             `json:"job_id"`      // Reference to the job with the anomaly
	Type        AnomalyType        `json:"type"`        // The specific type of anomaly
	Category    AnomalyCategory    `json:"category"`    // Whether it's absolute or relative
	Description string             `json:"description"` // Human-readable description
	Value       float64            `json:"value"`       // The actual value that triggered the anomaly
	Threshold   float64            `json:"threshold"`   // The threshold that was exceeded
	Operator    ComparisonOperator `json:"operator"`    // The comparison operator used
	CreatedAt   time.Time          `json:"created_at"`  // When the anomaly was detected
	Violations  []string           `json:"violations"`  // List of violations that led to the anomaly
	Severity    string             `json:"severity"`    // Severity of the anomaly
}

// Statistics holds statistical measures used for relative anomaly detection
type Statistics struct {
	// Salary statistics
	AvgSalary    float64
	SalaryStdDev float64

	// Requirements statistics
	AvgRequirements float64
	ReqStdDev       float64

	// Company rating statistics
	AvgRating    float64
	RatingStdDev float64

	// Location statistics
	AvgLatitude     float64
	LatitudeStdDev  float64
	AvgLongitude    float64
	LongitudeStdDev float64
}

// AnomalyService handles anomaly detection logic
type AnomalyService struct {
	db          DatabaseServiceInterface
	ruleService AnomalyRuleServiceInterface // Inject rule service for getting rules
}

// NewAnomalyService creates a new AnomalyService
func NewAnomalyService(db DatabaseServiceInterface, ruleService AnomalyRuleServiceInterface) *AnomalyService {
	return &AnomalyService{
		db:          db,
		ruleService: ruleService,
	}
}

// DetectAnomalies processes job data to detect anomalies based on rules
func (s *AnomalyService) DetectAnomalies(job *models.JobData) ([]models.Anomaly, error) {
	var detectedAnomalies []models.Anomaly

	// Check for null values in required fields
	var nullViolations []string
	if job.CompanyName == "" {
		nullViolations = append(nullViolations, "company_name")
	}
	if job.JobTitle == "" {
		nullViolations = append(nullViolations, "job_title")
	}
	if job.JobDescription == "" {
		nullViolations = append(nullViolations, "job_description")
	}
	if job.City == "" {
		nullViolations = append(nullViolations, "city")
	}
	if job.CompanyAddress == "" {
		nullViolations = append(nullViolations, "company_address")
	}
	if job.CompanyWebsite == "" {
		nullViolations = append(nullViolations, "company_website")
	}
	if job.JobLink == "" {
		nullViolations = append(nullViolations, "job_link")
	}

	// If there are null violations, create an anomaly
	if len(nullViolations) > 0 {
		nullAnomaly := models.Anomaly{
			Type:        models.AnomalyTypeNullValues,
			JobID:       job.JobID,
			Description: "Required fields are null",
			Value:       0,
			Threshold:   0,
			Operator:    models.Equal,
			CreatedAt:   time.Now(),
			Violations:  nullViolations,
		}
		if err := s.saveAnomaly(&nullAnomaly); err != nil {
			fmt.Printf("Error saving null value anomaly for job %s: %v\n", job.JobID, err)
		} else {
			detectedAnomalies = append(detectedAnomalies, nullAnomaly)
		}
	}

	// Get statistics for standard deviation checks
	stats, err := s.getStatistics()
	if err != nil {
		return nil, fmt.Errorf("error getting statistics: %w", err)
	}

	// Check for standard deviation anomalies in numeric fields
	if job.MaxSalary != nil {
		zScore := (*job.MaxSalary - stats.AvgSalary) / stats.SalaryStdDev
		if math.Abs(zScore) > StdDevThreshold {
			deviationAnomaly := models.Anomaly{
				Type:        models.AnomalyTypeDeviation,
				JobID:       job.JobID,
				Description: fmt.Sprintf("Salary deviates significantly from mean (z-score: %.2f)", zScore),
				Value:       *job.MaxSalary,
				Threshold:   stats.AvgSalary,
				Operator:    models.Equal,
				CreatedAt:   time.Now(),
				Violations:  []string{"max_salary"},
			}
			if err := s.saveAnomaly(&deviationAnomaly); err != nil {
				fmt.Printf("Error saving salary deviation anomaly for job %s: %v\n", job.JobID, err)
			} else {
				detectedAnomalies = append(detectedAnomalies, deviationAnomaly)
			}
		}
	}

	if job.CompanyRating != 0 {
		zScore := (job.CompanyRating - stats.AvgRating) / stats.RatingStdDev
		if math.Abs(zScore) > StdDevThreshold {
			deviationAnomaly := models.Anomaly{
				Type:        models.AnomalyTypeDeviation,
				JobID:       job.JobID,
				Description: fmt.Sprintf("Company rating deviates significantly from mean (z-score: %.2f)", zScore),
				Value:       job.CompanyRating,
				Threshold:   stats.AvgRating,
				Operator:    models.Equal,
				CreatedAt:   time.Now(),
				Violations:  []string{"company_rating"},
			}
			if err := s.saveAnomaly(&deviationAnomaly); err != nil {
				fmt.Printf("Error saving rating deviation anomaly for job %s: %v\n", job.JobID, err)
			} else {
				detectedAnomalies = append(detectedAnomalies, deviationAnomaly)
			}
		}
	}

	// Get active rules from the rule service
	rules, err := s.ruleService.GetAnomalyRules()
	if err != nil {
		return nil, fmt.Errorf("error getting anomaly rules via service: %w", err)
	}

	// Apply each active rule
	for _, rule := range rules {
		if !rule.IsActive {
			continue // Skip inactive rules
		}

		anomalyDetected := false
		var actualValue float64

		// Check based on rule type
		switch rule.Type {
		case models.AnomalyTypeMaxSalary:
			if job.MaxSalary != nil {
				actualValue = *job.MaxSalary
				anomalyDetected = compareValues(actualValue, rule.Value, rule.Operator)
			}
		case models.AnomalyTypeMinSalary:
			if job.MinSalary != nil {
				actualValue = *job.MinSalary
				anomalyDetected = compareValues(actualValue, rule.Value, rule.Operator)
			}
		case models.AnomalyTypeRating:
			// Assuming CompanyRating is not a pointer and always present
			actualValue = job.CompanyRating
			anomalyDetected = compareValues(actualValue, rule.Value, rule.Operator)
		default:
			// Log or handle unknown rule type if necessary
			continue
		}

		if anomalyDetected {
			anomaly := models.Anomaly{
				Type:        rule.Type,
				JobID:       job.JobID,
				Description: rule.Description,
				Value:       actualValue,
				Threshold:   rule.Value,
				Operator:    rule.Operator,
				CreatedAt:   time.Now(),
			}

			// Save the detected anomaly immediately
			if err := s.saveAnomaly(&anomaly); err != nil {
				// Log the error but continue processing other rules/anomalies
				fmt.Printf("Error saving anomaly for job %s, rule %d: %v\n", job.JobID, rule.ID, err)
			} else {
				detectedAnomalies = append(detectedAnomalies, anomaly)
			}
		}
	}

	return detectedAnomalies, nil
}

// getStatistics retrieves statistical measures for anomaly detection
func (s *AnomalyService) getStatistics() (*Statistics, error) {
	query := `
		SELECT 
			AVG(max_salary) as avg_salary,
			STDDEV(max_salary) as salary_stddev,
			AVG(company_rating) as avg_rating,
			STDDEV(company_rating) as rating_stddev
		FROM jobs
		WHERE max_salary IS NOT NULL AND company_rating > 0
	`

	var stats Statistics
	err := s.db.QueryRow(query).Scan(
		&stats.AvgSalary,
		&stats.SalaryStdDev,
		&stats.AvgRating,
		&stats.RatingStdDev,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting statistics: %w", err)
	}

	return &stats, nil
}

// saveAnomaly saves a single anomaly using basic exec methods
func (s *AnomalyService) saveAnomaly(anomaly *models.Anomaly) error {
	query := `
		INSERT INTO anomalies (job_id, type, description, value, threshold, operator, created_at, violations)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	// Use QueryRow as we need the ID back
	err := s.db.QueryRow(
		query,
		anomaly.JobID,
		anomaly.Type,
		anomaly.Description,
		anomaly.Value,
		anomaly.Threshold,
		anomaly.Operator,
		anomaly.CreatedAt,
		pq.Array(anomaly.Violations),
	).Scan(&anomaly.ID)

	if err != nil {
		return fmt.Errorf("error inserting anomaly: %w", err)
	}
	return nil
}

// compareValues performs the comparison based on the operator
func compareValues(value, threshold float64, operator models.ComparisonOperator) bool {
	switch operator {
	case models.GreaterThan:
		return value > threshold
	case models.GreaterThanOrEqual:
		return value >= threshold
	case models.LessThan:
		return value < threshold
	case models.LessThanOrEqual:
		return value <= threshold
	case models.Equal:
		return value == threshold
	default:
		return false // Unknown operator
	}
}

// GetAnomaliesByJobID retrieves anomalies for a specific job using basic query methods
func (s *AnomalyService) GetAnomaliesByJobID(jobID string) ([]models.Anomaly, error) {
	query := `
		SELECT id, job_id, type, description, value, threshold, operator, created_at
		FROM anomalies
		WHERE job_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, jobID)
	if err != nil {
		return nil, fmt.Errorf("error querying anomalies by job ID: %w", err)
	}
	defer rows.Close()

	var anomalies []models.Anomaly
	for rows.Next() {
		var anomaly models.Anomaly
		err := rows.Scan(
			&anomaly.ID,
			&anomaly.JobID,
			&anomaly.Type,
			&anomaly.Description,
			&anomaly.Value,
			&anomaly.Threshold,
			&anomaly.Operator,
			&anomaly.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning anomaly: %w", err)
		}
		anomalies = append(anomalies, anomaly)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating anomalies: %w", err)
	}

	if len(anomalies) == 0 {
		// Return empty slice and no error if no anomalies found, consistent with GetAnomalyRules
		return []models.Anomaly{}, nil
	}

	return anomalies, nil
}

// GetAllAnomalies retrieves all anomalies using basic query methods
func (s *AnomalyService) GetAllAnomalies() ([]models.Anomaly, error) {
	query := `
		SELECT id, job_id, type, description, value, threshold, operator, created_at
		FROM anomalies
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying all anomalies: %w", err)
	}
	defer rows.Close()

	var anomalies []models.Anomaly
	for rows.Next() {
		var anomaly models.Anomaly
		err := rows.Scan(
			&anomaly.ID,
			&anomaly.JobID,
			&anomaly.Type,
			&anomaly.Description,
			&anomaly.Value,
			&anomaly.Threshold,
			&anomaly.Operator,
			&anomaly.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning anomaly: %w", err)
		}
		anomalies = append(anomalies, anomaly)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating anomalies: %w", err)
	}

	return anomalies, nil
}

// DetectAnomaliesForAllJobs processes all existing jobs to detect anomalies
func (s *AnomalyService) DetectAnomaliesForAllJobs() error {
	// Get all jobs
	query := `
		SELECT job_id, company_name, company_rating, job_title, min_salary, max_salary
		FROM jobs
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return fmt.Errorf("error querying jobs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var job models.JobData
		err := rows.Scan(
			&job.JobID,
			&job.CompanyName,
			&job.CompanyRating,
			&job.JobTitle,
			&job.MinSalary,
			&job.MaxSalary,
		)
		if err != nil {
			return fmt.Errorf("error scanning job: %w", err)
		}

		// Detect anomalies for this job
		_, err = s.DetectAnomalies(&job)
		if err != nil {
			// Log the error but continue processing other jobs
			fmt.Printf("Error detecting anomalies for job %s: %v\n", job.JobID, err)
		}
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating jobs: %w", err)
	}

	return nil
}

// Removed rule management methods (GetAnomalyRules, GetAnomalyRule, CreateAnomalyRule, etc.)
// as they belong to AnomalyRuleService
