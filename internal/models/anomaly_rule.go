package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// StringSlice is a custom type for handling string arrays in JSON and database
type StringSlice []string

// Value implements the driver.Valuer interface
func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = StringSlice{}
		return nil
	}
	return json.Unmarshal(value.([]byte), s)
}

// AdvancedAnomalyRule represents a rule for advanced anomaly detection
type AdvancedAnomalyRule struct {
	ID          int64       `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Description string      `json:"description" db:"description"`
	Algorithm   string      `json:"algorithm" db:"algorithm"`
	Parameters  []byte      `json:"parameters" db:"parameters"`     // JSON
	InputFields StringSlice `json:"input_fields" db:"input_fields"` // JSON array
	Severity    string      `json:"severity" db:"severity"`
	IsActive    bool        `json:"is_active" db:"is_active"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name for the AdvancedAnomalyRule model
func (AdvancedAnomalyRule) TableName() string {
	return "anomaly_rules"
}
