package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ainesh01/anomaly_detection/internal/models"
)

// AnomalyRuleServiceInterface defines the interface for anomaly rule operations
type AnomalyRuleServiceInterface interface {
	GetAnomalyRules() ([]models.AnomalyRule, error)
	GetAnomalyRule(id int64) (*models.AnomalyRule, error)
	CreateAnomalyRule(rule *models.AnomalyRule) error
	UpdateAnomalyRule(rule *models.AnomalyRule) error
	DeleteAnomalyRule(id int64) error
	ToggleAnomalyRule(id int64, isActive bool) error
}

// AnomalyRuleService handles business logic for anomaly rules
type AnomalyRuleService struct {
	db DatabaseServiceInterface
}

// NewAnomalyRuleService creates a new AnomalyRuleService
func NewAnomalyRuleService(db DatabaseServiceInterface) *AnomalyRuleService {
	return &AnomalyRuleService{
		db: db,
	}
}

// GetAnomalyRules retrieves all anomaly rules using basic query methods
func (s *AnomalyRuleService) GetAnomalyRules() ([]models.AnomalyRule, error) {
	query := `
		SELECT id, name, description, type, operator, value, is_active, created_at, updated_at
		FROM anomaly_rules
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying anomaly rules: %w", err)
	}
	defer rows.Close()

	var rules []models.AnomalyRule
	for rows.Next() {
		var rule models.AnomalyRule
		err := rows.Scan(
			&rule.ID,
			&rule.Name,
			&rule.Description,
			&rule.Type,
			&rule.Operator,
			&rule.Value,
			&rule.IsActive,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning anomaly rule: %w", err)
		}
		rules = append(rules, rule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating anomaly rules: %w", err)
	}

	return rules, nil
}

// GetAnomalyRule retrieves a specific anomaly rule using basic query methods
func (s *AnomalyRuleService) GetAnomalyRule(id int64) (*models.AnomalyRule, error) {
	query := `
		SELECT id, name, description, type, operator, value, is_active, created_at, updated_at
		FROM anomaly_rules
		WHERE id = $1
	`

	var rule models.AnomalyRule
	row := s.db.QueryRow(query, id)
	err := row.Scan(
		&rule.ID,
		&rule.Name,
		&rule.Description,
		&rule.Type,
		&rule.Operator,
		&rule.Value,
		&rule.IsActive,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("anomaly rule with ID %d not found", id) // More specific error
		}
		return nil, fmt.Errorf("error querying or scanning anomaly rule: %w", err)
	}

	return &rule, nil
}

// CreateAnomalyRule creates a new anomaly rule using basic exec methods
func (s *AnomalyRuleService) CreateAnomalyRule(rule *models.AnomalyRule) error {
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = rule.CreatedAt // Set UpdatedAt to CreatedAt on creation

	query := `
		INSERT INTO anomaly_rules (name, description, type, operator, value, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	// Use QueryRow because we need the returned ID
	err := s.db.QueryRow(
		query,
		rule.Name,
		rule.Description,
		rule.Type,
		rule.Operator,
		rule.Value,
		rule.IsActive,
		rule.CreatedAt,
		rule.UpdatedAt,
	).Scan(&rule.ID)

	if err != nil {
		return fmt.Errorf("error creating anomaly rule: %w", err)
	}

	return nil
}

// UpdateAnomalyRule updates an existing anomaly rule using basic exec methods
func (s *AnomalyRuleService) UpdateAnomalyRule(rule *models.AnomalyRule) error {
	rule.UpdatedAt = time.Now()

	query := `
		UPDATE anomaly_rules
		SET name = $1,
			description = $2,
			type = $3,
			operator = $4,
			value = $5,
			is_active = $6,
			updated_at = $7
		WHERE id = $8
	`

	result, err := s.db.Exec(
		query,
		rule.Name,
		rule.Description,
		rule.Type,
		rule.Operator,
		rule.Value,
		rule.IsActive,
		rule.UpdatedAt,
		rule.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating anomaly rule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Log this error but don't necessarily fail the operation
		fmt.Printf("Could not get rows affected after update: %v\n", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("anomaly rule with ID %d not found for update", rule.ID)
	}

	return nil
}

// DeleteAnomalyRule deletes an anomaly rule using basic exec methods
func (s *AnomalyRuleService) DeleteAnomalyRule(id int64) error {
	query := `DELETE FROM anomaly_rules WHERE id = $1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting anomaly rule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Could not get rows affected after delete: %v\n", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("anomaly rule with ID %d not found for deletion", id)
	}

	return nil
}

// ToggleAnomalyRule toggles the active state of an anomaly rule using basic exec methods
func (s *AnomalyRuleService) ToggleAnomalyRule(id int64, isActive bool) error {
	query := `
		UPDATE anomaly_rules
		SET is_active = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	result, err := s.db.Exec(query, isActive, id)
	if err != nil {
		return fmt.Errorf("error toggling anomaly rule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Could not get rows affected after toggle: %v\n", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("anomaly rule with ID %d not found for toggle", id)
	}

	return nil
}
