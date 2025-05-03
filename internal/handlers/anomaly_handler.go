package handlers

import (
	"net/http"

	"github.com/ainesh01/anomaly_detection/internal/models"
	"github.com/ainesh01/anomaly_detection/internal/services"
	"github.com/gin-gonic/gin"
)

// AnomalyHandler handles HTTP requests for anomalies
type AnomalyHandler struct {
	anomalyService services.AnomalyServiceInterface
}

// NewAnomalyHandler creates a new AnomalyHandler
func NewAnomalyHandler(anomalyService services.AnomalyServiceInterface) *AnomalyHandler {
	return &AnomalyHandler{
		anomalyService: anomalyService,
	}
}

// GetAnomaliesByJobID handles GET requests for anomalies by job ID
func (h *AnomalyHandler) GetAnomaliesByJobID(c *gin.Context) {
	jobID := c.Param("job_id")
	anomalies, err := h.anomalyService.GetAnomaliesByJobID(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, anomalies)
}

// GetAllAnomalies handles GET requests for all anomalies
func (h *AnomalyHandler) GetAllAnomalies(c *gin.Context) {
	anomalies, err := h.anomalyService.GetAllAnomalies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if anomalies == nil {
		anomalies = []models.Anomaly{} // Ensure we return an empty array instead of null
	}
	c.JSON(http.StatusOK, anomalies)
}

// DetectAnomalies handles POST request to detect anomalies for a job
func (h *AnomalyHandler) DetectAnomalies(c *gin.Context) {
	var jobData models.JobData
	if err := c.ShouldBindJSON(&jobData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	anomalies, err := h.anomalyService.DetectAnomalies(&jobData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, anomalies)
}

// DetectAnomaliesForAllJobs handles POST request to detect anomalies for all jobs
func (h *AnomalyHandler) DetectAnomaliesForAllJobs(c *gin.Context) {
	if err := h.anomalyService.DetectAnomaliesForAllJobs(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Anomaly detection completed for all jobs"})
}
