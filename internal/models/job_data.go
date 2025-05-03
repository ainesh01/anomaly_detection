package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// CustomTime is a wrapper around time.Time that implements custom JSON and SQL marshaling
type CustomTime struct {
	time.Time
}

// Value implements the driver.Valuer interface
func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time, nil
}

// Scan implements the sql.Scanner interface
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		ct.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		ct.Time = v
		return nil
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
		ct.Time = t
		return nil
	default:
		return fmt.Errorf("cannot scan %T into CustomTime", value)
	}
}

// MarshalJSON implements the json.Marshaler interface
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(ct.Time.Format(time.RFC3339))
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		ct.Time = time.Time{}
		return nil
	}

	// Try different time formats
	formats := []string{
		time.RFC3339,                  // "2006-01-02T15:04:05Z07:00"
		"2006-01-02 15:04:05.999 MST", // "2025-03-23 01:43:50.322 UTC"
		"2006-01-02 15:04:05 MST",     // "2025-03-23 01:43:50 UTC"
		"2006-01-02 15:04:05.999",     // "2025-03-23 01:43:50.322"
		"2006-01-02 15:04:05",         // "2025-03-23 01:43:50"
	}

	var lastErr error
	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			ct.Time = t
			return nil
		} else {
			lastErr = err
		}
	}

	return fmt.Errorf("could not parse time %q with any known format: %v", s, lastErr)
}

// JobData represents a job listing with all its associated data
type JobData struct {
	// Company Information
	CompanyName    string  `json:"companyName"`
	CompanyRating  float64 `json:"companyRating"`
	CompanyAddress string  `json:"companyAddress"`
	CompanyWebsite string  `json:"companyWebsite"`

	// Job Information
	JobTitle         string     `json:"jobTitle"`
	JobPostedTime    CustomTime `json:"jobPostedTime"`
	JobID            string     `json:"jobID"`
	JobLink          string     `json:"jobLink"`
	JobDescription   string     `json:"jobDescription"`
	JobRequirements  []string   `json:"jobRequirements"`
	JobBenefits      []string   `json:"jobBenefits"`
	JobTypes         []string   `json:"jobTypes"`
	IsNewJob         bool       `json:"isNewJob"`
	IsNoResumeJob    bool       `json:"isNoResumeJob"`
	IsUrgentlyHiring bool       `json:"isUrgentlyHiring"`

	// Role Information
	RoleType          *string  `json:"roleType,omitempty"`
	MinSalary         *float64 `json:"minSalary,omitempty"`
	MaxSalary         *float64 `json:"maxSalary,omitempty"`
	SalaryGranularity *string  `json:"salaryGranularity,omitempty"`
	HiresNeeded       *string  `json:"hiresNeeded,omitempty"`

	// Location Information
	City          string   `json:"city"`
	State         *string  `json:"state,omitempty"`
	Zip           *string  `json:"zip,omitempty"`
	PlaceID       *string  `json:"placeId,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	LocationCount int      `json:"locationCount"`

	// Social Media Links
	Facebook  *string `json:"facebook,omitempty"`
	Instagram *string `json:"instagram,omitempty"`
	Tiktok    *string `json:"tiktok,omitempty"`
	Youtube   *string `json:"youtube,omitempty"`
	Twitter   *string `json:"twitter,omitempty"`
	Yelp      *string `json:"yelp,omitempty"`

	// Additional Information
	SchedulingLink *string `json:"schedulingLink,omitempty"`

	// Metadata
	InvocationID    string     `json:"invocationID"`
	TaskID          string     `json:"taskID"`
	DateRepresented CustomTime `json:"dateRepresented"`
	DateCollected   CustomTime `json:"dateCollected"`
	AttemptID       string     `json:"attemptID"`

	// Database timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
