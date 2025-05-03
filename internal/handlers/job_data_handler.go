package handlers

import (
	"net/http"

	"github.com/ainesh01/anomaly_detection/internal/models"
	"github.com/ainesh01/anomaly_detection/internal/services"
	"github.com/gin-gonic/gin"
)

// JobDataHandler handles HTTP requests for job data
type JobDataHandler struct {
	jobDataService services.JobDataServiceInterface
}

// NewJobDataHandler creates a new JobDataHandler
func NewJobDataHandler(jobDataService services.JobDataServiceInterface) *JobDataHandler {
	return &JobDataHandler{
		jobDataService: jobDataService,
	}
}

// CreateJobData handles POST requests to create a new job data entry
func (h *JobDataHandler) CreateJobData(c *gin.Context) {
	var job models.JobData
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.jobDataService.CreateJobData(&job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, job)
}

// GetJobData handles GET requests for a specific job data entry
func (h *JobDataHandler) GetJobData(c *gin.Context) {
	jobID := c.Param("job_id")
	job, err := h.jobDataService.GetJobData(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}

// GetAllJobData handles GET requests for all job data entries
func (h *JobDataHandler) GetAllJobData(c *gin.Context) {
	jobs, err := h.jobDataService.GetAllJobData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobs)
}
