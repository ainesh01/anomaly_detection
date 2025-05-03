package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ainesh01/anomaly_detection/internal/models"
	"github.com/lib/pq" // Needed for pq.Array
)

// JobDataServiceInterface defines the interface for job data service operations
type JobDataServiceInterface interface {
	CreateJobData(job *models.JobData) error
	GetJobData(jobID string) (*models.JobData, error)
	GetAllJobData() ([]models.JobData, error)
}

// JobDataService handles business logic for job data operations
type JobDataService struct {
	db DatabaseServiceInterface
}

// NewJobDataService creates a new JobDataService
func NewJobDataService(db DatabaseServiceInterface) *JobDataService {
	return &JobDataService{
		db: db,
	}
}

// CreateJobData creates or updates a job data entry using basic exec methods
func (s *JobDataService) CreateJobData(job *models.JobData) error {
	// Set timestamps
	now := time.Now()
	if job.CreatedAt.IsZero() {
		job.CreatedAt = now
	}
	job.UpdatedAt = now

	// Use ON CONFLICT to handle potential existing job_id
	query := `
		INSERT INTO jobs (
			job_id, company_name, company_rating, company_address, company_website,
			job_title, job_posted_time, job_link, job_description,
			job_requirements, job_benefits, job_types, is_new_job,
			is_no_resume_job, is_urgently_hiring, role_type, min_salary,
			max_salary, salary_granularity, hires_needed, city, state,
			zip, place_id, latitude, longitude, location_count, facebook,
			instagram, tiktok, youtube, twitter, yelp, scheduling_link,
			invocation_id, task_id, date_represented, date_collected, attempt_id,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26,
			$27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41
		)
		ON CONFLICT (job_id) DO UPDATE SET
			company_name = EXCLUDED.company_name,
			company_rating = EXCLUDED.company_rating,
			company_address = EXCLUDED.company_address,
			company_website = EXCLUDED.company_website,
			job_title = EXCLUDED.job_title,
			job_posted_time = EXCLUDED.job_posted_time,
			job_link = EXCLUDED.job_link,
			job_description = EXCLUDED.job_description,
			job_requirements = EXCLUDED.job_requirements,
			job_benefits = EXCLUDED.job_benefits,
			job_types = EXCLUDED.job_types,
			is_new_job = EXCLUDED.is_new_job,
			is_no_resume_job = EXCLUDED.is_no_resume_job,
			is_urgently_hiring = EXCLUDED.is_urgently_hiring,
			role_type = EXCLUDED.role_type,
			min_salary = EXCLUDED.min_salary,
			max_salary = EXCLUDED.max_salary,
			salary_granularity = EXCLUDED.salary_granularity,
			hires_needed = EXCLUDED.hires_needed,
			city = EXCLUDED.city,
			state = EXCLUDED.state,
			zip = EXCLUDED.zip,
			place_id = EXCLUDED.place_id,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			location_count = EXCLUDED.location_count,
			facebook = EXCLUDED.facebook,
			instagram = EXCLUDED.instagram,
			tiktok = EXCLUDED.tiktok,
			youtube = EXCLUDED.youtube,
			twitter = EXCLUDED.twitter,
			yelp = EXCLUDED.yelp,
			scheduling_link = EXCLUDED.scheduling_link,
			invocation_id = EXCLUDED.invocation_id,
			task_id = EXCLUDED.task_id,
			date_represented = EXCLUDED.date_represented,
			date_collected = EXCLUDED.date_collected,
			attempt_id = EXCLUDED.attempt_id,
			updated_at = EXCLUDED.updated_at
	`

	_, err := s.db.Exec(query,
		job.JobID,
		job.CompanyName,
		job.CompanyRating,
		job.CompanyAddress,
		job.CompanyWebsite,
		job.JobTitle,
		job.JobPostedTime,
		job.JobLink,
		job.JobDescription,
		pq.Array(job.JobRequirements),
		pq.Array(job.JobBenefits),
		pq.Array(job.JobTypes),
		job.IsNewJob,
		job.IsNoResumeJob,
		job.IsUrgentlyHiring,
		job.RoleType,
		job.MinSalary,
		job.MaxSalary,
		job.SalaryGranularity,
		job.HiresNeeded,
		job.City,
		job.State,
		job.Zip,
		job.PlaceID,
		job.Latitude,
		job.Longitude,
		job.LocationCount,
		job.Facebook,
		job.Instagram,
		job.Tiktok,
		job.Youtube,
		job.Twitter,
		job.Yelp,
		job.SchedulingLink,
		job.InvocationID,
		job.TaskID,
		job.DateRepresented,
		job.DateCollected,
		job.AttemptID,
		job.CreatedAt,
		job.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error saving job data: %w", err)
	}

	return nil
}

// GetJobData retrieves a specific job data entry using basic query methods
func (s *JobDataService) GetJobData(jobID string) (*models.JobData, error) {
	// Select all columns from the jobs table
	query := `
		SELECT
			job_id, company_name, company_rating, company_address, company_website,
			job_title, job_posted_time, job_link, job_description,
			job_requirements, job_benefits, job_types, is_new_job,
			is_no_resume_job, is_urgently_hiring, role_type, min_salary,
			max_salary, salary_granularity, hires_needed, city, state,
			zip, place_id, latitude, longitude, location_count, facebook,
			instagram, tiktok, youtube, twitter, yelp, scheduling_link,
			invocation_id, task_id, date_represented, date_collected, attempt_id,
			created_at, updated_at
		FROM jobs
		WHERE job_id = $1
	`

	row := s.db.QueryRow(query, jobID)
	job := &models.JobData{}

	// Scan all columns into the JobData struct
	err := row.Scan(
		&job.JobID,
		&job.CompanyName,
		&job.CompanyRating,
		&job.CompanyAddress,
		&job.CompanyWebsite,
		&job.JobTitle,
		&job.JobPostedTime,
		&job.JobLink,
		&job.JobDescription,
		pq.Array(&job.JobRequirements),
		pq.Array(&job.JobBenefits),
		pq.Array(&job.JobTypes),
		&job.IsNewJob,
		&job.IsNoResumeJob,
		&job.IsUrgentlyHiring,
		&job.RoleType,
		&job.MinSalary,
		&job.MaxSalary,
		&job.SalaryGranularity,
		&job.HiresNeeded,
		&job.City,
		&job.State,
		&job.Zip,
		&job.PlaceID,
		&job.Latitude,
		&job.Longitude,
		&job.LocationCount,
		&job.Facebook,
		&job.Instagram,
		&job.Tiktok,
		&job.Youtube,
		&job.Twitter,
		&job.Yelp,
		&job.SchedulingLink,
		&job.InvocationID,
		&job.TaskID,
		&job.DateRepresented,
		&job.DateCollected,
		&job.AttemptID,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job data with ID %s not found", jobID)
		}
		return nil, fmt.Errorf("error querying or scanning job data: %w", err)
	}

	return job, nil
}

// GetAllJobData retrieves all job data entries
func (s *JobDataService) GetAllJobData() ([]models.JobData, error) {
	// Select all fields from the jobs table
	query := `
		SELECT
			job_id, company_name, company_rating, company_address, company_website,
			job_title, job_posted_time, job_link, job_description,
			job_requirements, job_benefits, job_types, is_new_job,
			is_no_resume_job, is_urgently_hiring, role_type, min_salary,
			max_salary, salary_granularity, hires_needed, city, state,
			zip, place_id, latitude, longitude, location_count, facebook,
			instagram, tiktok, youtube, twitter, yelp, scheduling_link,
			invocation_id, task_id, date_represented, date_collected, attempt_id,
			created_at, updated_at
		FROM jobs
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying all job data: %w", err)
	}
	defer rows.Close()

	var jobs []models.JobData
	for rows.Next() {
		var job models.JobData
		// Scan all fields into the JobData struct
		err := rows.Scan(
			&job.JobID,
			&job.CompanyName,
			&job.CompanyRating,
			&job.CompanyAddress,
			&job.CompanyWebsite,
			&job.JobTitle,
			&job.JobPostedTime,
			&job.JobLink,
			&job.JobDescription,
			pq.Array(&job.JobRequirements),
			pq.Array(&job.JobBenefits),
			pq.Array(&job.JobTypes),
			&job.IsNewJob,
			&job.IsNoResumeJob,
			&job.IsUrgentlyHiring,
			&job.RoleType,
			&job.MinSalary,
			&job.MaxSalary,
			&job.SalaryGranularity,
			&job.HiresNeeded,
			&job.City,
			&job.State,
			&job.Zip,
			&job.PlaceID,
			&job.Latitude,
			&job.Longitude,
			&job.LocationCount,
			&job.Facebook,
			&job.Instagram,
			&job.Tiktok,
			&job.Youtube,
			&job.Twitter,
			&job.Yelp,
			&job.SchedulingLink,
			&job.InvocationID,
			&job.TaskID,
			&job.DateRepresented,
			&job.DateCollected,
			&job.AttemptID,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning job data row: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating job data rows: %w", err)
	}

	return jobs, nil
}
