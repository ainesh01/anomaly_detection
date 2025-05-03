package models

import "time"

// AnomalyRuleExecution represents the execution of an anomaly detection rule
type AnomalyRuleExecution struct {
	ID          int64      `json:"id" db:"id"`
	RuleID      int64      `json:"rule_id" db:"rule_id"`
	Status      string     `json:"status" db:"status"`
	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	Result      []byte     `json:"result,omitempty" db:"result"` // JSON
	Error       *string    `json:"error,omitempty" db:"error"`
}
