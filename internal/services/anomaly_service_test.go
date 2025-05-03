package services

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/ainesh01/anomaly_detection/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDatabaseService is a mock implementation of DatabaseServiceInterface
type MockDatabaseService struct {
	mock.Mock
}

func (m *MockDatabaseService) Exec(query string, args ...interface{}) (sql.Result, error) {
	arguments := m.Called(query, args)
	return arguments.Get(0).(sql.Result), arguments.Error(1)
}

func (m *MockDatabaseService) Query(query string, args ...interface{}) (*sql.Rows, error) {
	arguments := m.Called(query, args)
	return arguments.Get(0).(*sql.Rows), arguments.Error(1)
}

func (m *MockDatabaseService) QueryRow(query string, args ...interface{}) *sql.Row {
	arguments := m.Called(query, args)
	return arguments.Get(0).(*sql.Row)
}

func (m *MockDatabaseService) Close() error {
	arguments := m.Called()
	return arguments.Error(0)
}

// MockAnomalyRuleService is a mock implementation of AnomalyRuleServiceInterface
type MockAnomalyRuleService struct {
	mock.Mock
}

func (m *MockAnomalyRuleService) GetAnomalyRules() ([]models.AnomalyRule, error) {
	arguments := m.Called()
	return arguments.Get(0).([]models.AnomalyRule), arguments.Error(1)
}

func (m *MockAnomalyRuleService) GetAnomalyRule(id int64) (*models.AnomalyRule, error) {
	args := m.Called(id)
	return args.Get(0).(*models.AnomalyRule), args.Error(1)
}

func (m *MockAnomalyRuleService) CreateAnomalyRule(rule *models.AnomalyRule) error {
	args := m.Called(rule)
	return args.Error(0)
}

func (m *MockAnomalyRuleService) UpdateAnomalyRule(rule *models.AnomalyRule) error {
	args := m.Called(rule)
	return args.Error(0)
}

func (m *MockAnomalyRuleService) DeleteAnomalyRule(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAnomalyRuleService) ToggleAnomalyRule(id int64, isActive bool) error {
	args := m.Called(id, isActive)
	return args.Error(0)
}

func TestDetectAnomalies(t *testing.T) {
	// Create mock services
	mockDB := new(MockDatabaseService)
	mockRuleService := new(MockAnomalyRuleService)

	// Create anomaly service with mocks
	service := NewAnomalyService(mockDB, mockRuleService)

	t.Run("Test Null Value Detection", func(t *testing.T) {
		// Create a job with missing required fields
		job := &models.JobData{
			JobID: "test-job-1",
			// CompanyName is empty
			JobTitle:       "Software Engineer",
			JobDescription: "Test description",
			City:           "Test City",
			// CompanyAddress is empty
			// CompanyWebsite is empty
			// JobLink is empty
		}

		// Set up mock expectations for null value check
		mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(&sql.Row{})
		mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil)

		// Call DetectAnomalies
		anomalies, err := service.DetectAnomalies(job)
		assert.NoError(t, err)
		assert.NotEmpty(t, anomalies)

		// Verify null value anomalies
		var nullValueAnomaly *models.Anomaly
		for _, anomaly := range anomalies {
			if anomaly.Type == models.AnomalyTypeNullValues {
				nullValueAnomaly = &anomaly
				break
			}
		}
		assert.NotNil(t, nullValueAnomaly)
		assert.Contains(t, nullValueAnomaly.Violations, "company_name")
		assert.Contains(t, nullValueAnomaly.Violations, "company_address")
		assert.Contains(t, nullValueAnomaly.Violations, "company_website")
		assert.Contains(t, nullValueAnomaly.Violations, "job_link")
	})

	t.Run("Test Standard Deviation Detection", func(t *testing.T) {
		// Create a job with extreme values
		maxSalary := 500000.0 // Very high salary
		companyRating := 5.0  // Perfect rating
		job := &models.JobData{
			JobID:          "test-job-2",
			CompanyName:    "Test Company",
			JobTitle:       "Software Engineer",
			JobDescription: "Test description",
			City:           "Test City",
			CompanyAddress: "Test Address",
			CompanyWebsite: "test.com",
			JobLink:        "test.com/job",
			MaxSalary:      &maxSalary,
			CompanyRating:  companyRating,
		}

		// Set up mock expectations for statistics query
		statsRow := &sql.Row{}
		mockDB.On("QueryRow", mock.Anything).Return(statsRow)
		mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil)

		// Mock statistics data
		stats := &Statistics{
			AvgSalary:    100000.0,
			SalaryStdDev: 50000.0,
			AvgRating:    3.5,
			RatingStdDev: 0.5,
		}

		// Calculate expected z-scores
		expectedSalaryZScore := (maxSalary - stats.AvgSalary) / stats.SalaryStdDev
		expectedRatingZScore := (companyRating - stats.AvgRating) / stats.RatingStdDev

		// Call DetectAnomalies
		anomalies, err := service.DetectAnomalies(job)
		assert.NoError(t, err)
		assert.NotEmpty(t, anomalies)

		// Verify standard deviation anomalies
		var salaryDeviationAnomaly *models.Anomaly
		var ratingDeviationAnomaly *models.Anomaly
		for _, anomaly := range anomalies {
			if anomaly.Type == models.AnomalyTypeDeviation {
				if len(anomaly.Violations) == 1 && anomaly.Violations[0] == "max_salary" {
					salaryDeviationAnomaly = &anomaly
				} else if len(anomaly.Violations) == 1 && anomaly.Violations[0] == "company_rating" {
					ratingDeviationAnomaly = &anomaly
				}
			}
		}
		assert.NotNil(t, salaryDeviationAnomaly)
		assert.NotNil(t, ratingDeviationAnomaly)

		// Verify z-scores in descriptions
		assert.Contains(t, salaryDeviationAnomaly.Description, fmt.Sprintf("z-score: %.2f", expectedSalaryZScore))
		assert.Contains(t, ratingDeviationAnomaly.Description, fmt.Sprintf("z-score: %.2f", expectedRatingZScore))
	})

	t.Run("Test Rule-Based Detection", func(t *testing.T) {
		// Create a job with values that should trigger rules
		maxSalary := -1000.0 // Negative salary
		job := &models.JobData{
			JobID:          "test-job-3",
			CompanyName:    "Test Company",
			JobTitle:       "Software Engineer",
			JobDescription: "Test description",
			City:           "Test City",
			CompanyAddress: "Test Address",
			CompanyWebsite: "test.com",
			JobLink:        "test.com/job",
			MaxSalary:      &maxSalary,
			CompanyRating:  3.5,
		}

		// Set up mock expectations for rules
		rules := []models.AnomalyRule{
			{
				ID:          1,
				Name:        "Negative Salary",
				Description: "Alert if maximum salary is negative",
				Type:        models.AnomalyTypeSalary,
				Operator:    models.LessThan,
				Value:       0.0,
				IsActive:    true,
			},
		}
		mockRuleService.On("GetAnomalyRules").Return(rules, nil)
		mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(&sql.Row{})
		mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil)

		// Call DetectAnomalies
		anomalies, err := service.DetectAnomalies(job)
		assert.NoError(t, err)
		assert.NotEmpty(t, anomalies)

		// Verify rule-based anomaly
		var ruleAnomaly *models.Anomaly
		for _, anomaly := range anomalies {
			if anomaly.Type == models.AnomalyTypeSalary {
				ruleAnomaly = &anomaly
				break
			}
		}
		assert.NotNil(t, ruleAnomaly)
		assert.Equal(t, models.AnomalyTypeSalary, ruleAnomaly.Type)
		assert.Equal(t, maxSalary, ruleAnomaly.Value)
		assert.Equal(t, 0.0, ruleAnomaly.Threshold)
		assert.Equal(t, models.LessThan, ruleAnomaly.Operator)
	})
}
