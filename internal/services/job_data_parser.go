package services

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/ainesh01/anomaly_detection/internal/models"
)

// ParseJSONLFile reads a JSONL file (optionally gzipped) and returns a slice of JobData
func ParseJSONLFile(filePath string) ([]models.JobData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var reader *bufio.Reader

	// Check if the file is gzipped
	if strings.HasSuffix(filepath.Base(filePath), ".gz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
		defer gzReader.Close()
		reader = bufio.NewReader(gzReader)
	} else {
		reader = bufio.NewReader(file)
	}

	var jobs []models.JobData
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var job models.JobData
		if err := json.Unmarshal(scanner.Bytes(), &job); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}
