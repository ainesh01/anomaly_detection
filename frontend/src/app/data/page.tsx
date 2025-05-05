'use client';

import { useState, useEffect } from 'react';
import { JobData } from '@/types'; // Assuming src is mapped to @/
import { columns } from "./columns"; // Import columns
import { DataTable } from "@/components/ui/data-table"; // Assuming a reusable DataTable component exists
import { Input } from "@/components/ui/input"; // Import Input component

export default function DataPage() {
  const [jobs, setJobs] = useState<JobData[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState<string>(''); // State for search query

  useEffect(() => {
    async function fetchData() {
      try {
        // Assuming the backend runs on port 8080 locally
        const response = await fetch('http://localhost:8080/api/job-data');
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data: JobData[] = await response.json();
        setJobs(data);
      } catch (e: any) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    }

    fetchData();
  }, []); // Empty dependency array ensures this runs only once on mount

  // Filter jobs based on search query
  const filteredJobs = jobs.filter(job =>
    job.jobID.toLowerCase().includes(searchQuery.toLowerCase())
  );

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error fetching data: {error}</div>;

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Job Data</h1>
      {/* Add search input */}
      <div className="mb-4">
        <Input
          placeholder="Search by Job ID..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="max-w-sm"
        />
      </div>
      {/* Pass filteredJobs to DataTable */}
      <DataTable columns={columns} data={filteredJobs} />
      {/* <pre className="bg-gray-100 p-4 rounded overflow-auto text-sm">
        {JSON.stringify(jobs, null, 2)}
      </pre> */}
    </div>
  );
} 