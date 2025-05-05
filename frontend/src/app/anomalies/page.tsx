"use client"; // Make this a Client Component

import * as React from "react"; // Import React, useState, useEffect
import { useState } from "react";
import { DataTable } from "@/components/ui/data-table";
import { Button } from "@/components/ui/button"; // Import Button
import type { ColumnDef } from "@tanstack/react-table";

// Define a basic type for the Anomaly data (adjust based on actual API response)
type Anomaly = {
  id: string;
  timestamp: string; // Or Date?
  rule_name: string; // Name of the rule that triggered this
  metric: string;
  value: number;
  threshold: number;
  severity: string;
};

// We move the fetching logic inside the component using useEffect
// async function getAnomalies(): Promise<Anomaly[]> { ... }

export default function AnomaliesPage() {
  const [anomaliesData, setAnomaliesData] = useState<Anomaly[]>([]);
  const [isLoading, setIsLoading] = useState(false); // Initially not loading
  const [error, setError] = useState<string | null>(null);

  // Function to handle running the anomaly check
  async function runAnomalyCheck() {
    setIsLoading(true);
    setError(null);
    setAnomaliesData([]); // Clear previous results
    try {
      // Use the correct endpoint found: POST /api/anomalies/detect-all
      const response = await fetch('http://localhost:8080/api/anomalies/detect-all', { // Corrected path
          method: 'POST',
          // Add headers or body if required by the endpoint
          // headers: { 'Content-Type': 'application/json' },
          // body: JSON.stringify({ /* potential parameters */ }),
      });
      if (!response.ok) {
         // Improved error handling for different statuses
         let errorMsg = `HTTP error! status: ${response.status}`;
         try {
           const errorBody = await response.json(); // Try to parse backend error message
           if (errorBody && errorBody.error) {
             errorMsg += `: ${errorBody.error}`;
           }
         } catch (parseError) {
           // Ignore if body isn't JSON or empty
         }
         throw new Error(errorMsg);
      }

      // After detection, fetch the new anomalies list
      // Assuming GET /api/anomalies returns the list
      const anomaliesResponse = await fetch('http://localhost:8080/api/anomalies');
      if (!anomaliesResponse.ok) {
         throw new Error(`HTTP error fetching anomalies! status: ${anomaliesResponse.status}`);
      }
      const data = await anomaliesResponse.json();
      setAnomaliesData(data);

    } catch (e) {
       if (e instanceof Error) {
          setError(`Failed to run anomaly check: ${e.message}`);
      } else {
          setError("An unknown error occurred while running anomaly check");
      }
      console.error("Fetch error:", e);
    } finally {
      setIsLoading(false);
    }
  }

  // Update the columns definition to show only Job ID and Description
  const columns: ColumnDef<Anomaly>[] = [
    {
        accessorKey: "job_id",
        header: "Job ID",
      },
      {
        accessorKey: "description",
        header: "Description",
      },
      // Remove other columns
    // {
    //     accessorKey: "timestamp", // Keep or remove based on the API response field name
    //     header: "Timestamp",
    //     cell: ({ row }) => {
    //         const date = new Date(row.getValue("created_at")); // Assuming 'created_at' is the correct field
    //         return date.toLocaleString();
    //     }
    //   },
    //   {
    //     accessorKey: "rule_name", // Check if 'type' should be used instead
    //     header: "Rule Name / Type",
    //   },
    //   {
    //     accessorKey: "metric", // This seems missing from the example response
    //     header: "Metric",
    //   },
    //   {
    //     accessorKey: "value",
    //     header: "Detected Value",
    //   },
    //   {
    //       accessorKey: "threshold",
    //       header: "Threshold",
    //   },
    //   {
    //     accessorKey: "severity", // This seems missing from the example response
    //     header: "Severity",
    //   },
  ];

  return (
    <div className="container mx-auto py-10">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">Run Anomaly Detection</h1>
        <Button onClick={runAnomalyCheck} disabled={isLoading}>
          {isLoading ? "Running..." : "Run Check"}
        </Button>
      </div>

      {/* Display loading state */}
      {isLoading && <div className="text-center py-4">Running anomaly detection...</div>}

      {/* Display error message */}
      {error && <div className="text-center py-4 text-red-500">Error: {error}</div>}

      {/* Display DataTable only when not loading, no error, and data exists */}
      {!isLoading && !error && (
         <DataTable columns={columns} data={anomaliesData} />
      )}
    </div>
  );
} 