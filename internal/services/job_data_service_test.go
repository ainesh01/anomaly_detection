package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/ainesh01/anomaly_detection/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestJobDataService(t *testing.T) {
	t.Run("CreateJobData", func(t *testing.T) {
		// Setup
		mockDB := new(MockDB)
		service := NewJobDataService(mockDB)
		job := &models.JobData{
			JobID:          "job1",
			JobTitle:       "Software Engineer",
			CompanyName:    "Tech Corp",
			JobDescription: "Job description",
			MinSalary:      Float64Ptr(50000.0),
			MaxSalary:      Float64Ptr(100000.0),
			JobRequirements: []string{
				"Go",
				"Python",
			},
			CompanyRating:   4.5,
			Latitude:        Float64Ptr(37.7749),
			Longitude:       Float64Ptr(-122.4194),
			JobPostedTime:   models.CustomTime{Time: time.Now()},
			DateRepresented: models.CustomTime{Time: time.Now()},
			DateCollected:   models.CustomTime{Time: time.Now()},
		}

		// Setup mock result
		mockResult := new(MockResult)
		mockResult.On("RowsAffected").Return(int64(1), nil)
		mockDB.On("Exec", mock.Anything, mock.Anything).Return(mockResult, nil)

		// Test
		err := service.CreateJobData(job)

		// Assertions
		assert.NoError(t, err)
		assert.False(t, job.CreatedAt.IsZero())
		assert.False(t, job.UpdatedAt.IsZero())
		mockDB.AssertExpectations(t)
	})

	t.Run("GetJobData", func(t *testing.T) {
		// Setup
		mockDB := new(MockDB)
		service := NewJobDataService(mockDB)
		expectedJob := &models.JobData{
			JobID:          "job1",
			JobTitle:       "Software Engineer",
			CompanyName:    "Tech Corp",
			JobDescription: "Job description",
			MinSalary:      Float64Ptr(50000.0),
			MaxSalary:      Float64Ptr(100000.0),
			JobRequirements: []string{
				"Go",
				"Python",
			},
			CompanyRating:   4.5,
			Latitude:        Float64Ptr(37.7749),
			Longitude:       Float64Ptr(-122.4194),
			JobPostedTime:   models.CustomTime{Time: time.Now()},
			DateRepresented: models.CustomTime{Time: time.Now()},
			DateCollected:   models.CustomTime{Time: time.Now()},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		// Setup mock row
		mockRow := &sql.Row{}
		mockDB.On("QueryRow", mock.Anything, "job1").Return(mockRow)

		// Test
		job, err := service.GetJobData("job1")

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedJob, job)
		mockDB.AssertExpectations(t)
	})

	t.Run("GetAllJobData", func(t *testing.T) {
		// Setup
		mockDB := new(MockDB)
		service := NewJobDataService(mockDB)
		expectedJobs := []models.JobData{
			{
				JobID:          "job1",
				JobTitle:       "Software Engineer",
				CompanyName:    "Tech Corp",
				JobDescription: "Job description",
				MinSalary:      Float64Ptr(50000.0),
				MaxSalary:      Float64Ptr(100000.0),
				JobRequirements: []string{
					"Go",
					"Python",
				},
				CompanyRating:   4.5,
				Latitude:        Float64Ptr(37.7749),
				Longitude:       Float64Ptr(-122.4194),
				JobPostedTime:   models.CustomTime{Time: time.Now()},
				DateRepresented: models.CustomTime{Time: time.Now()},
				DateCollected:   models.CustomTime{Time: time.Now()},
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				JobID:          "job2",
				JobTitle:       "Data Scientist",
				CompanyName:    "Data Corp",
				JobDescription: "Data job description",
				MinSalary:      Float64Ptr(60000.0),
				MaxSalary:      Float64Ptr(120000.0),
				JobRequirements: []string{
					"Python",
					"R",
				},
				CompanyRating:   4.0,
				Latitude:        Float64Ptr(37.7749),
				Longitude:       Float64Ptr(-122.4194),
				JobPostedTime:   models.CustomTime{Time: time.Now()},
				DateRepresented: models.CustomTime{Time: time.Now()},
				DateCollected:   models.CustomTime{Time: time.Now()},
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		}

		// Setup mock rows
		mockRows := &sql.Rows{}
		mockDB.On("Query", mock.Anything).Return(mockRows, nil)

		// Test
		jobs, err := service.GetAllJobData()

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedJobs, jobs)
		mockDB.AssertExpectations(t)
	})

	t.Run("Error Cases", func(t *testing.T) {
		// Setup
		mockDB := new(MockDB)
		service := NewJobDataService(mockDB)
		expectedError := assert.AnError

		t.Run("CreateJobData Error", func(t *testing.T) {
			job := &models.JobData{
				JobID:          "job1",
				JobTitle:       "Software Engineer",
				CompanyName:    "Tech Corp",
				JobDescription: "Job description",
				MinSalary:      Float64Ptr(50000.0),
				MaxSalary:      Float64Ptr(100000.0),
				JobRequirements: []string{
					"Go",
					"Python",
				},
				CompanyRating:   4.5,
				Latitude:        Float64Ptr(37.7749),
				Longitude:       Float64Ptr(-122.4194),
				JobPostedTime:   models.CustomTime{Time: time.Now()},
				DateRepresented: models.CustomTime{Time: time.Now()},
				DateCollected:   models.CustomTime{Time: time.Now()},
			}
			mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, expectedError)
			err := service.CreateJobData(job)
			assert.Error(t, err)
			assert.Equal(t, expectedError, err)
		})

		t.Run("GetJobData Error", func(t *testing.T) {
			mockDB.On("QueryRow", mock.Anything, "job1").Return(nil)
			job, err := service.GetJobData("job1")
			assert.Error(t, err)
			assert.Nil(t, job)
		})

		t.Run("GetAllJobData Error", func(t *testing.T) {
			mockDB.On("Query", mock.Anything).Return(nil, expectedError)
			jobs, err := service.GetAllJobData()
			assert.Error(t, err)
			assert.Nil(t, jobs)
			assert.Equal(t, expectedError, err)
		})
	})
}
