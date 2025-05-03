package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ainesh01/anomaly_detection/internal/config"
	"github.com/ainesh01/anomaly_detection/internal/handlers"
	"github.com/ainesh01/anomaly_detection/internal/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	servercfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	dbcfg := config.NewDBConfig()

	// Initialize database service
	dbService, err := services.InitializeDatabaseService(dbcfg)
	if err != nil {
		log.Fatalf("Error initializing database service: %v", err)
	}
	defer dbService.Close()

	// Initialize services
	jobDataService := services.NewJobDataService(dbService)
	anomalyRuleService := services.NewAnomalyRuleService(dbService)
	anomalyService := services.NewAnomalyService(dbService, anomalyRuleService)

	// Check if a file was provided
	filePath := parseCommandLineArgs()
	if filePath != "" {
		// Parse the file and detect anomalies
		rows, err := services.ParseJSONLFile(filePath)
		if err != nil {
			log.Fatalf("Error parsing file: %v", err)
		}

		// Save each job to the database
		for _, job := range rows {
			if err := jobDataService.CreateJobData(&job); err != nil {
				log.Printf("Error saving job %s: %v", job.JobID, err)
				continue
			}
		}
		log.Printf("Successfully parsed and saved %d rows from %s", len(rows), filePath)
	} else {
		log.Fatal("No file provided. Please provide a file to parse.")
	}

	// Initialize HTTP server
	srv := setupServer(jobDataService, anomalyService, anomalyRuleService, servercfg)

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

// parseCommandLineArgs parses and validates command line arguments
// Returns the file path to parse or empty string if not provided
func parseCommandLineArgs() string {
	filePath := flag.String("file", "", "Path to the JSONL.gz file to parse")
	flag.Parse()
	return *filePath
}

func setupServer(
	jobDataService services.JobDataServiceInterface,
	anomalyService services.AnomalyServiceInterface,
	anomalyRuleService services.AnomalyRuleServiceInterface,
	servercfg *config.ServerConfig,
) *http.Server {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Initialize handlers
	jobDataHandler := handlers.NewJobDataHandler(jobDataService)
	anomalyHandler := handlers.NewAnomalyHandler(anomalyService)
	anomalyRuleHandler := handlers.NewAnomalyRuleHandler(anomalyRuleService)

	// Define API endpoints
	api := router.Group("/api")
	{
		// Job data endpoints
		api.POST("/job-data", jobDataHandler.CreateJobData)
		api.GET("/job-data/:job_id", jobDataHandler.GetJobData)
		api.GET("/job-data", jobDataHandler.GetAllJobData)

		// Anomaly endpoints
		api.GET("/anomalies/:job_id", anomalyHandler.GetAnomaliesByJobID)
		api.GET("/anomalies", anomalyHandler.GetAllAnomalies)
		api.POST("/anomalies/detect-all", anomalyHandler.DetectAnomaliesForAllJobs)

		// Anomaly rule endpoints
		api.GET("/anomaly-rules", anomalyRuleHandler.GetAnomalyRules)
		api.GET("/anomaly-rules/:id", anomalyRuleHandler.GetAnomalyRule)
		api.POST("/anomaly-rules", anomalyRuleHandler.CreateAnomalyRule)
		api.PUT("/anomaly-rules/:id", anomalyRuleHandler.UpdateAnomalyRule)
		api.DELETE("/anomaly-rules/:id", anomalyRuleHandler.DeleteAnomalyRule)
		api.PATCH("/anomaly-rules/:id/toggle", anomalyRuleHandler.ToggleAnomalyRule)
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", servercfg.Port),
		Handler: router,
	}
}
