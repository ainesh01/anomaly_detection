package models

import "time"

// AnomalyAlert represents an alert generated when an anomaly is detected
type AnomalyAlert struct {
	ID          int64      `json:"id" db:"id"`
	RuleID      int64      `json:"rule_id" db:"rule_id"`
	Severity    string     `json:"severity" db:"severity"`
	Description string     `json:"description" db:"description"`
	Details     []byte     `json:"details" db:"details"` // JSON
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty" db:"resolved_at"`
	Status      string     `json:"status" db:"status"`
}
