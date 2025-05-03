package handlers

import (
	"net/http"
	"strconv"

	"github.com/ainesh01/anomaly_detection/internal/models"
	"github.com/ainesh01/anomaly_detection/internal/services"
	"github.com/gin-gonic/gin"
)

// AnomalyRuleHandler handles HTTP requests for anomaly rules
type AnomalyRuleHandler struct {
	ruleService services.AnomalyRuleServiceInterface
}

// NewAnomalyRuleHandler creates a new AnomalyRuleHandler
func NewAnomalyRuleHandler(ruleService services.AnomalyRuleServiceInterface) *AnomalyRuleHandler {
	return &AnomalyRuleHandler{
		ruleService: ruleService,
	}
}

// GetAnomalyRules handles GET requests for all anomaly rules
func (h *AnomalyRuleHandler) GetAnomalyRules(c *gin.Context) {
	rules, err := h.ruleService.GetAnomalyRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}

// GetAnomalyRule handles GET requests for a specific anomaly rule
func (h *AnomalyRuleHandler) GetAnomalyRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	rule, err := h.ruleService.GetAnomalyRule(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

// CreateAnomalyRule handles POST requests to create a new anomaly rule
func (h *AnomalyRuleHandler) CreateAnomalyRule(c *gin.Context) {
	var rule models.AnomalyRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ruleService.CreateAnomalyRule(&rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rule)
}

// UpdateAnomalyRule handles PUT requests to update an existing anomaly rule
func (h *AnomalyRuleHandler) UpdateAnomalyRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	var rule models.AnomalyRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule.ID = id
	if err := h.ruleService.UpdateAnomalyRule(&rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

// DeleteAnomalyRule handles DELETE requests to remove an anomaly rule
func (h *AnomalyRuleHandler) DeleteAnomalyRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	if err := h.ruleService.DeleteAnomalyRule(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ToggleAnomalyRule handles PATCH requests to toggle the active state of an anomaly rule
func (h *AnomalyRuleHandler) ToggleAnomalyRule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	var request struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ruleService.ToggleAnomalyRule(id, request.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
