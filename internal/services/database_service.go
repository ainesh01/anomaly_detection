package services

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ainesh01/anomaly_detection/internal/config"
)

// DatabaseServiceInterface defines the interface for basic database operations
type DatabaseServiceInterface interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}

// SQLDB is a concrete implementation of DatabaseServiceInterface using *sql.DB
type SQLDB struct {
	db *sql.DB
}

// InitializeDatabaseService sets up the database connection and creates tables.
// Returns the simplified DatabaseServiceInterface.
func InitializeDatabaseService(cfg *config.DBConfig) (DatabaseServiceInterface, error) {
	dbService, err := NewDatabaseService(cfg) // This now returns DatabaseServiceInterface (SQLDB)
	if err != nil {
		log.Fatalf("Error initializing database service: %v", err) // Use Fatalf for critical init errors
	}
	// Keep defer dbService.Close() in main.go where the service is used

	// Create database tables using the interface
	if err := createTables(dbService); err != nil {
		dbService.Close() // Attempt to close before fatal exit
		log.Fatalf("Error creating tables: %v", err)
	}

	return dbService, nil
}

// NewDatabaseService creates a new database connection wrapped by SQLDB.
// Returns the simplified DatabaseServiceInterface.
func NewDatabaseService(cfg *config.DBConfig) (DatabaseServiceInterface, error) {
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close() // Close if ping fails
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Database connection successful")
	return &SQLDB{db: db}, nil
}

// Exec executes a query without returning rows.
func (s *SQLDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

// Query executes a query that returns rows.
func (s *SQLDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
func (s *SQLDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.db.QueryRow(query, args...)
}

// Close closes the database connection.
func (s *SQLDB) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// createTables creates the necessary database tables if they don't exist.
// It now accepts the interface to execute queries.
func createTables(dbService DatabaseServiceInterface) error {
	// Drop tables in reverse order of dependencies
	dropQueries := []string{
		`DROP TABLE IF EXISTS anomalies;`,
		`DROP TABLE IF EXISTS jobs;`,
		`DROP TABLE IF EXISTS anomaly_rules;`,
	}

	for _, query := range dropQueries {
		_, err := dbService.Exec(query)
		if err != nil {
			return fmt.Errorf("error dropping tables: %v", err)
		}
	}

	// Create tables in order of dependencies
	if err := createJobsTable(dbService); err != nil {
		return err
	}
	if err := createAnomaliesTable(dbService); err != nil {
		return err
	}
	if err := createAnomalyRulesTable(dbService); err != nil {
		return err
	}

	// Create default anomaly rules
	if err := createDefaultAnomalyRules(dbService); err != nil {
		return err
	}

	return nil
}

func createJobsTable(dbService DatabaseServiceInterface) error {
	query := `
		CREATE TABLE jobs (
			job_id TEXT PRIMARY KEY,
			company_name TEXT NOT NULL,
			company_rating DOUBLE PRECISION,
			company_address TEXT,
			company_website TEXT,
			job_title TEXT NOT NULL,
			job_posted_time TIMESTAMP WITH TIME ZONE,
			job_link TEXT,
			job_description TEXT,
			job_requirements TEXT[],
			job_benefits TEXT[],
			job_types TEXT[],
			is_new_job BOOLEAN,
			is_no_resume_job BOOLEAN,
			is_urgently_hiring BOOLEAN,
			role_type TEXT,
			min_salary DOUBLE PRECISION,
			max_salary DOUBLE PRECISION,
			salary_granularity TEXT,
			hires_needed TEXT,
			city TEXT,
			state TEXT,
			zip TEXT,
			place_id TEXT,
			latitude DOUBLE PRECISION,
			longitude DOUBLE PRECISION,
			location_count INTEGER,
			facebook TEXT,
			instagram TEXT,
			tiktok TEXT,
			youtube TEXT,
			twitter TEXT,
			yelp TEXT,
			scheduling_link TEXT,
			invocation_id TEXT,
			task_id TEXT,
			date_represented TIMESTAMP WITH TIME ZONE,
			date_collected TIMESTAMP WITH TIME ZONE,
			attempt_id TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := dbService.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating jobs table: %v", err)
	}
	log.Println("Jobs table created successfully.")
	return nil
}

// Added anomalies table creation based on model fields previously used
func createAnomaliesTable(dbService DatabaseServiceInterface) error {
	query := `
		CREATE TABLE anomalies (
			id BIGSERIAL PRIMARY KEY,
			job_id TEXT NOT NULL REFERENCES jobs(job_id),
			type TEXT NOT NULL,
			description TEXT NOT NULL,
			value DOUBLE PRECISION,
			threshold DOUBLE PRECISION,
			operator TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			violations TEXT[]
		);

		CREATE INDEX idx_anomalies_job_id ON anomalies(job_id);
		CREATE INDEX idx_anomalies_type ON anomalies(type);
	`
	_, err := dbService.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating anomalies table: %v", err)
	}
	log.Println("Anomalies table created successfully.")
	return nil
}

func createAnomalyRulesTable(dbService DatabaseServiceInterface) error {
	query := `
		CREATE TABLE anomaly_rules (
			id BIGSERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			description TEXT NOT NULL,
			type TEXT NOT NULL,
			operator TEXT NOT NULL,
			value DOUBLE PRECISION NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX idx_anomaly_rules_name ON anomaly_rules(name);
		CREATE INDEX idx_anomaly_rules_active ON anomaly_rules(is_active);
	`

	_, err := dbService.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating anomaly rules table: %v", err)
	}
	log.Println("Anomaly rules table created successfully.")
	return nil
}

// createDefaultAnomalyRules creates some default rules for anomaly detection
func createDefaultAnomalyRules(dbService DatabaseServiceInterface) error {
	query := `
		INSERT INTO anomaly_rules (name, description, type, operator, value, is_active, created_at, updated_at)
		VALUES 
		('Negative Salary', 'Alert if maximum salary is negative', 'salary', '<', 0.0, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (name) DO NOTHING;
	`

	_, err := dbService.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating default anomaly rules: %v", err)
	}
	log.Println("Default anomaly rules created successfully.")
	return nil
}

// Removed createAnomalyRuleExecutionsTable and createAnomalyAlertsTable
// as they seemed related to the more complex AdvancedAnomalyRule setup.

// Removed specific data methods like InsertJob, GetJobByID, GetJobByRowIndex,
// GetJobsByRowIndexRange, GetAllJobs, SaveAnomalyRuleExecution, SaveAnomalyAlert,
// GetActiveAnomalyRules, GetAnomalyRules, GetAnomalyRule, SaveAnomalyRule,
// ToggleAnomalyRule, DeleteAnomalyRule, SaveAnomaly, GetAnomaliesByJobID,
// GetAllAnomalies, GetJobData, SaveJobData.
// This logic will be moved to the specific service implementations (JobDataService, AnomalyService, etc.)
// where they will use the dbService.Exec, dbService.Query, dbService.QueryRow methods.
