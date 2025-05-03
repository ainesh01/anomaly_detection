package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ainesh01/anomaly_detection/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRuleDB is a mock implementation of DatabaseServiceInterface
type MockRuleDB struct {
	mock.Mock
}

func (m *MockRuleDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	arguments := m.Called(query, args)
	return arguments.Get(0).(sql.Result), arguments.Error(1)
}

func (m *MockRuleDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	arguments := m.Called(query, args)
	return arguments.Get(0).(*sql.Rows), arguments.Error(1)
}

func (m *MockRuleDB) QueryRow(query string, args ...interface{}) *sql.Row {
	arguments := m.Called(query, args)
	return arguments.Get(0).(*sql.Row)
}

func (m *MockRuleDB) Close() error {
	arguments := m.Called()
	return arguments.Error(0)
}

func TestAnomalyRuleService(t *testing.T) {
	t.Run("GetAnomalyRules", func(t *testing.T) {
		// Create SQL mock
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Create expected rows
		columns := []string{"id", "name", "description", "type", "operator", "value", "is_active", "created_at", "updated_at"}
		mock.ExpectQuery("SELECT (.+) FROM anomaly_rules").WillReturnRows(
			sqlmock.NewRows(columns).AddRow(
				1,
				"Test Rule",
				"Test Description",
				models.AnomalyTypeSalary,
				models.GreaterThan,
				100000.0,
				true,
				time.Now(),
				time.Now(),
			),
		)

		// Create mock DB that returns the sqlmock rows
		mockDB := new(MockRuleDB)
		expectedQuery := "\n\t\tSELECT id, name, description, type, operator, value, is_active, created_at, updated_at\n\t\tFROM anomaly_rules\n\t\tORDER BY created_at DESC\n\t"
		mockDB.On("Query", expectedQuery, []interface{}(nil)).Return(db.Query("SELECT * FROM anomaly_rules"))

		// Create service with mock
		service := NewAnomalyRuleService(mockDB)

		// Call GetAnomalyRules
		rules, err := service.GetAnomalyRules()
		assert.NoError(t, err)
		assert.NotEmpty(t, rules)
		assert.Equal(t, 1, len(rules))
		assert.Equal(t, "Test Rule", rules[0].Name)

		// Verify mock expectations
		mockDB.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAnomalyRule", func(t *testing.T) {
		// Create SQL mock
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Create mock DB
		mockDB := new(MockRuleDB)

		// Create service with mock
		service := NewAnomalyRuleService(mockDB)

		// Create expected rows
		now := time.Now()
		mock.ExpectQuery("SELECT (.+) FROM anomaly_rules").WithArgs(1).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "description", "type", "operator", "value", "is_active", "created_at", "updated_at"}).
				AddRow(1, "Test Rule", "Test Description", models.AnomalyTypeSalary, models.GreaterThan, 100000.0, true, now, now),
		)

		// Set up mock expectations
		expectedQuery := "\n\t\tSELECT id, name, description, type, operator, value, is_active, created_at, updated_at\n\t\tFROM anomaly_rules\n\t\tWHERE id = $1\n\t"
		mockDB.On("QueryRow", expectedQuery, []interface{}{int64(1)}).Return(db.QueryRow("SELECT * FROM anomaly_rules WHERE id = $1", 1))

		// Call GetAnomalyRule
		rule, err := service.GetAnomalyRule(1)
		assert.NoError(t, err)
		assert.NotNil(t, rule)
		assert.Equal(t, "Test Rule", rule.Name)
		assert.Equal(t, models.AnomalyTypeSalary, rule.Type)
		assert.Equal(t, models.GreaterThan, rule.Operator)
		assert.Equal(t, 100000.0, rule.Value)
		assert.True(t, rule.IsActive)

		// Verify mock expectations
		mockDB.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateAnomalyRule", func(t *testing.T) {
		// Create SQL mock
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Create mock DB
		mockDB := new(MockRuleDB)

		// Create service with mock
		service := NewAnomalyRuleService(mockDB)

		// Use a fixed timestamp for testing
		fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		// Create test rule
		rule := &models.AnomalyRule{
			Name:        "Test Rule",
			Description: "Test Description",
			Type:        models.AnomalyTypeSalary,
			Operator:    models.GreaterThan,
			Value:       100000.0,
			IsActive:    true,
		}

		// Set up mock expectations for insert
		mock.ExpectQuery("INSERT INTO anomaly_rules").WithArgs(
			rule.Name,
			rule.Description,
			rule.Type,
			rule.Operator,
			rule.Value,
			rule.IsActive,
			sqlmock.AnyArg(), // created_at can be any time
			sqlmock.AnyArg(), // updated_at can be any time
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Set up mock expectations
		expectedQuery := "\n\t\tINSERT INTO anomaly_rules (name, description, type, operator, value, is_active, created_at, updated_at)\n\t\tVALUES ($1, $2, $3, $4, $5, $6, $7, $8)\n\t\tRETURNING id\n\t"
		mockDB.On("QueryRow", expectedQuery, []interface{}{
			rule.Name,
			rule.Description,
			rule.Type,
			rule.Operator,
			rule.Value,
			rule.IsActive,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		}).Return(db.QueryRow("INSERT INTO anomaly_rules (name, description, type, operator, value, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
			rule.Name,
			rule.Description,
			rule.Type,
			rule.Operator,
			rule.Value,
			rule.IsActive,
			fixedTime,
			fixedTime,
		))

		// Call CreateAnomalyRule
		err = service.CreateAnomalyRule(rule)
		assert.NoError(t, err)

		// Verify mock expectations
		mockDB.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())

		// Verify the timestamps were set
		assert.False(t, rule.CreatedAt.IsZero())
		assert.False(t, rule.UpdatedAt.IsZero())
		assert.Equal(t, rule.CreatedAt, rule.UpdatedAt)
	})

	t.Run("UpdateAnomalyRule", func(t *testing.T) {
		// Setup
		mockDB := new(MockRuleDB)
		service := NewAnomalyRuleService(mockDB)
		rule := &models.AnomalyRule{
			ID:          1,
			Name:        "High Salary Check",
			Description: "Alert if salary exceeds $200,000",
			Type:        models.AnomalyTypeSalary,
			Operator:    models.GreaterThan,
			Value:       200000.0,
			IsActive:    true,
			CreatedAt:   time.Now(),
		}

		// Setup mock result
		mockResult := new(MockResult)
		mockResult.On("RowsAffected").Return(int64(1), nil)
		mockDB.On("Exec", mock.Anything, mock.Anything).Return(mockResult, nil)

		// Test
		err := service.UpdateAnomalyRule(rule)

		// Assertions
		assert.NoError(t, err)
		assert.False(t, rule.UpdatedAt.IsZero())
		mockDB.AssertExpectations(t)
	})

	t.Run("DeleteAnomalyRule", func(t *testing.T) {
		// Setup
		mockDB := new(MockRuleDB)
		service := NewAnomalyRuleService(mockDB)

		// Setup mock result
		mockResult := new(MockResult)
		mockResult.On("RowsAffected").Return(int64(1), nil)
		mockDB.On("Exec", mock.Anything, int64(1)).Return(mockResult, nil)

		// Test
		err := service.DeleteAnomalyRule(1)

		// Assertions
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("ToggleAnomalyRule", func(t *testing.T) {
		// Create SQL mock
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Create mock DB
		mockDB := new(MockRuleDB)

		// Create service with mock
		service := NewAnomalyRuleService(mockDB)

		// Set up mock expectations
		mock.ExpectExec("UPDATE anomaly_rules").WillReturnResult(sqlmock.NewResult(1, 1))
		expectedQuery := "\n\t\tUPDATE anomaly_rules\n\t\tSET is_active = $1\n\t\tWHERE id = $2\n\t"
		mockDB.On("Exec", expectedQuery, []interface{}{true, int64(1)}).Return(db.Exec("UPDATE anomaly_rules"))

		// Call ToggleAnomalyRule
		err = service.ToggleAnomalyRule(1, true)
		assert.NoError(t, err)

		// Verify mock expectations
		mockDB.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error Cases", func(t *testing.T) {
		// Setup
		mockDB := new(MockRuleDB)
		service := NewAnomalyRuleService(mockDB)
		expectedError := assert.AnError

		t.Run("GetAnomalyRules Error", func(t *testing.T) {
			mockDB.On("Query", mock.Anything).Return(nil, expectedError)
			rules, err := service.GetAnomalyRules()
			assert.Error(t, err)
			assert.Nil(t, rules)
			assert.Equal(t, expectedError, err)
		})

		t.Run("GetAnomalyRule Error", func(t *testing.T) {
			mockDB.On("QueryRow", mock.Anything, int64(1)).Return(nil)
			rule, err := service.GetAnomalyRule(1)
			assert.Error(t, err)
			assert.Nil(t, rule)
		})

		t.Run("CreateAnomalyRule Error", func(t *testing.T) {
			rule := &models.AnomalyRule{
				Name:        "High Salary Check",
				Description: "Alert if salary exceeds $200,000",
				Type:        models.AnomalyTypeSalary,
				Operator:    models.GreaterThan,
				Value:       200000.0,
				IsActive:    true,
			}
			mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, expectedError)
			err := service.CreateAnomalyRule(rule)
			assert.Error(t, err)
			assert.Equal(t, expectedError, err)
		})

		t.Run("UpdateAnomalyRule Error", func(t *testing.T) {
			rule := &models.AnomalyRule{
				ID:          1,
				Name:        "High Salary Check",
				Description: "Alert if salary exceeds $200,000",
				Type:        models.AnomalyTypeSalary,
				Operator:    models.GreaterThan,
				Value:       200000.0,
				IsActive:    true,
				CreatedAt:   time.Now(),
			}
			mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, expectedError)
			err := service.UpdateAnomalyRule(rule)
			assert.Error(t, err)
			assert.Equal(t, expectedError, err)
		})

		t.Run("DeleteAnomalyRule Error", func(t *testing.T) {
			mockDB.On("Exec", mock.Anything, int64(1)).Return(nil, expectedError)
			err := service.DeleteAnomalyRule(1)
			assert.Error(t, err)
			assert.Equal(t, expectedError, err)
		})

		t.Run("ToggleAnomalyRule Error", func(t *testing.T) {
			mockDB.On("Exec", mock.Anything, int64(1), false).Return(nil, expectedError)
			err := service.ToggleAnomalyRule(1, false)
			assert.Error(t, err)
			assert.Equal(t, expectedError, err)
		})
	})
}
