package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/ainesh01/anomaly_detection/internal/config"
	"github.com/ainesh01/anomaly_detection/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of DatabaseServiceInterface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	arguments := m.Called(query, args)
	return arguments.Get(0).(sql.Result), arguments.Error(1)
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	arguments := m.Called(query, args)
	return arguments.Get(0).(*sql.Rows), arguments.Error(1)
}

func (m *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	arguments := m.Called(query, args)
	return arguments.Get(0).(*sql.Row)
}

func (m *MockDB) Close() error {
	arguments := m.Called()
	return arguments.Error(0)
}

func (m *MockDB) GetJobsByRowIndexRange(start, end int64) ([]models.JobData, error) {
	arguments := m.Called(start, end)
	return arguments.Get(0).([]models.JobData), arguments.Error(1)
}

func (m *MockDB) GetAllJobs() ([]models.JobData, error) {
	arguments := m.Called()
	return arguments.Get(0).([]models.JobData), arguments.Error(1)
}

func (m *MockDB) SaveJobData(job *models.JobData) error {
	arguments := m.Called(job)
	return arguments.Error(0)
}

func (m *MockDB) GetJobData(jobID string) (*models.JobData, error) {
	arguments := m.Called(jobID)
	return arguments.Get(0).(*models.JobData), arguments.Error(1)
}

func (m *MockDB) SaveAnomaly(anomaly *Anomaly) error {
	arguments := m.Called(anomaly)
	return arguments.Error(0)
}

func (m *MockDB) GetAnomaliesByJobID(jobID string) ([]Anomaly, error) {
	arguments := m.Called(jobID)
	return arguments.Get(0).([]Anomaly), arguments.Error(1)
}

func (m *MockDB) GetAllAnomalies() ([]Anomaly, error) {
	arguments := m.Called()
	return arguments.Get(0).([]Anomaly), arguments.Error(1)
}

func (m *MockDB) SaveAnomalyRule(rule *models.AdvancedAnomalyRule) error {
	arguments := m.Called(rule)
	return arguments.Error(0)
}

func (m *MockDB) GetAnomalyRule(id int64) (*models.AdvancedAnomalyRule, error) {
	arguments := m.Called(id)
	return arguments.Get(0).(*models.AdvancedAnomalyRule), arguments.Error(1)
}

func (m *MockDB) GetActiveAnomalyRules() ([]*models.AdvancedAnomalyRule, error) {
	arguments := m.Called()
	return arguments.Get(0).([]*models.AdvancedAnomalyRule), arguments.Error(1)
}

func (m *MockDB) ToggleAnomalyRule(id int64, isActive bool) error {
	arguments := m.Called(id, isActive)
	return arguments.Error(0)
}

// MockResult is a mock implementation of sql.Result
type MockResult struct {
	mock.Mock
}

func (m *MockResult) LastInsertId() (int64, error) {
	arguments := m.Called()
	return arguments.Get(0).(int64), arguments.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	arguments := m.Called()
	return arguments.Get(0).(int64), arguments.Error(1)
}

func TestInitializeDatabaseService(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.DBConfig
		expectError bool
	}{
		{
			name: "valid configuration",
			config: &config.DBConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "postgres",
				DBName:   "anomaly_detection_test",
			},
			expectError: false,
		},
		{
			name: "invalid configuration",
			config: &config.DBConfig{
				Host:     "invalid_host",
				Port:     5432,
				User:     "postgres",
				Password: "postgres",
				DBName:   "anomaly_detection_test",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := InitializeDatabaseService(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				if service != nil {
					defer service.Close()
				}
			}
		})
	}
}

func TestGetJobsByRowIndexRange(t *testing.T) {
	// Create mock database
	mockDB := new(MockDB)

	// Test data
	now := time.Now()
	testJobs := []models.JobData{
		{
			JobID:           "test1",
			JobTitle:        "Software Engineer",
			CompanyName:     "Test Company 1",
			JobDescription:  "Test job description 1",
			JobPostedTime:   models.CustomTime{Time: now},
			DateRepresented: models.CustomTime{Time: now},
			DateCollected:   models.CustomTime{Time: now},
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		{
			JobID:           "test2",
			JobTitle:        "Data Scientist",
			CompanyName:     "Test Company 2",
			JobDescription:  "Test job description 2",
			JobPostedTime:   models.CustomTime{Time: now},
			DateRepresented: models.CustomTime{Time: now},
			DateCollected:   models.CustomTime{Time: now},
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		{
			JobID:           "test3",
			JobTitle:        "Product Manager",
			CompanyName:     "Test Company 3",
			JobDescription:  "Test job description 3",
			JobPostedTime:   models.CustomTime{Time: now},
			DateRepresented: models.CustomTime{Time: now},
			DateCollected:   models.CustomTime{Time: now},
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}

	tests := []struct {
		name         string
		start        int64
		end          int64
		expectedJobs []models.JobData
		expectError  bool
		setupMock    func()
	}{
		{
			name:         "valid range",
			start:        1,
			end:          3,
			expectedJobs: testJobs,
			expectError:  false,
			setupMock: func() {
				// Setup mock to return test jobs
				mockDB.On("GetJobsByRowIndexRange", int64(1), int64(3)).
					Return(testJobs, nil)
			},
		},
		{
			name:         "empty range",
			start:        10,
			end:          20,
			expectedJobs: []models.JobData{},
			expectError:  false,
			setupMock: func() {
				// Setup mock to return empty result
				mockDB.On("GetJobsByRowIndexRange", int64(10), int64(20)).
					Return([]models.JobData{}, nil)
			},
		},
		{
			name:         "invalid range",
			start:        5,
			end:          3,
			expectedJobs: nil,
			expectError:  true,
			setupMock:    func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			tt.setupMock()

			// Call the function
			jobs, err := mockDB.GetJobsByRowIndexRange(tt.start, tt.end)

			// Assert results
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, jobs)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, jobs)
				assert.Equal(t, len(tt.expectedJobs), len(jobs))
			}

			// Verify mock expectations
			mockDB.AssertExpectations(t)
		})
	}
}

func TestGetAllJobs(t *testing.T) {
	// Create mock database
	mockDB := new(MockDB)

	// Test data
	now := time.Now()
	testJobs := []models.JobData{
		{
			JobID:           "test1",
			JobTitle:        "Software Engineer",
			CompanyName:     "Test Company 1",
			JobDescription:  "Test job description 1",
			JobPostedTime:   models.CustomTime{Time: now},
			DateRepresented: models.CustomTime{Time: now},
			DateCollected:   models.CustomTime{Time: now},
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		{
			JobID:           "test2",
			JobTitle:        "Data Scientist",
			CompanyName:     "Test Company 2",
			JobDescription:  "Test job description 2",
			JobPostedTime:   models.CustomTime{Time: now},
			DateRepresented: models.CustomTime{Time: now},
			DateCollected:   models.CustomTime{Time: now},
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		{
			JobID:           "test3",
			JobTitle:        "Product Manager",
			CompanyName:     "Test Company 3",
			JobDescription:  "Test job description 3",
			JobPostedTime:   models.CustomTime{Time: now},
			DateRepresented: models.CustomTime{Time: now},
			DateCollected:   models.CustomTime{Time: now},
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}

	// Setup mock
	mockDB.On("GetAllJobs").Return(testJobs, nil)

	// Call the function
	jobs, err := mockDB.GetAllJobs()

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, jobs)
	assert.Equal(t, len(testJobs), len(jobs))

	// Verify mock expectations
	mockDB.AssertExpectations(t)
}
