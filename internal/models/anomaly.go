package models

import "time"

type AnomalyType string
type ComparisonOperator string

const (
	// Simple predefined check types
	AnomalyTypeMaxSalary  AnomalyType = "max_salary"         // For max salary threshold checks
	AnomalyTypeMinSalary  AnomalyType = "min_salary"         // For min salary threshold checks
	AnomalyTypeRating     AnomalyType = "company_rating"     // For company rating checks
	AnomalyTypeNullValues AnomalyType = "null_values"        // For null value checks
	AnomalyTypeDeviation  AnomalyType = "standard_deviation" // For standard deviation checks

	// Operators
	GreaterThan        ComparisonOperator = ">"
	GreaterThanOrEqual ComparisonOperator = ">="
	LessThan           ComparisonOperator = "<"
	LessThanOrEqual    ComparisonOperator = "<="
	Equal              ComparisonOperator = "="
)

// Anomaly represents a detected anomaly
type Anomaly struct {
	ID          string             `json:"id"`
	Type        AnomalyType        `json:"type"`
	JobID       string             `json:"job_id"`
	Description string             `json:"description"`
	Value       float64            `json:"value"`
	Threshold   float64            `json:"threshold"`
	Operator    ComparisonOperator `json:"operator"`
	CreatedAt   time.Time          `json:"created_at"`
	Violations  []string           `json:"violations"` // List of fields that violated the rule
}

// AnomalyRule represents a simple predefined check rule
type AnomalyRule struct {
	ID          int64              `json:"id" db:"id"`
	Name        string             `json:"name" db:"name"`
	Description string             `json:"description" db:"description"`
	Type        AnomalyType        `json:"type" db:"type"`           // Type of check (salary, rating)
	Operator    ComparisonOperator `json:"operator" db:"operator"`   // The comparison operator
	Value       float64            `json:"value" db:"value"`         // The threshold value
	IsActive    bool               `json:"is_active" db:"is_active"` // Whether the rule is active
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name for the AnomalyRule model
func (AnomalyRule) TableName() string {
	return "anomaly_rules"
}

// AnomalyRuleRequest represents the data needed to create or update a rule
type AnomalyRuleRequest struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description" binding:"required"`
	Type        AnomalyType        `json:"type" binding:"required"`
	Operator    ComparisonOperator `json:"operator" binding:"required"`
	Value       float64            `json:"value" binding:"required"`
	IsActive    bool               `json:"is_active"`
}
