"use client"

import { ColumnDef } from "@tanstack/react-table"
import { JobData } from "@/types"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"

// Helper function to render optional fields
const renderOptional = (value: string | number | null | undefined, label: string) => {
  if (value === null || value === undefined || value === "") {
    return null
  }
  return <p><strong>{label}:</strong> {value}</p>
}

// Helper function to render array fields
const renderArray = (value: string[] | undefined, label: string) => {
  if (!value || value.length === 0) {
    return null
  }
  return <p><strong>{label}:</strong> {value.join(", ")}</p>
}

export const columns: ColumnDef<JobData>[] = [
  {
    accessorKey: "jobID",
    header: "Job ID",
  },
  {
    accessorKey: "jobTitle",
    header: "Job Title",
  },
  {
    accessorKey: "companyName",
    header: "Company Name",
  },
  {
    id: "actions",
    cell: ({ row }: { row: any }) => {
      const job = row.original as JobData

      return (
        <Dialog>
          <DialogTrigger asChild>
            <Button variant="outline">View Details</Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[600px] max-h-[80vh] overflow-y-auto overflow-x-hidden">
            <DialogHeader>
              <DialogTitle>{job.jobTitle}</DialogTitle>
              <DialogDescription>
                {job.companyName} - {job.city}{job.state ? `, ${job.state}` : ""}
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4 text-sm">
              {/* Display all fields from JobData */}
              <p><strong>Job ID:</strong> {job.jobID}</p>
              <p><strong>Company Rating:</strong> {job.companyRating}</p>
              <p><strong>Company Address:</strong> {job.companyAddress}</p>
              <p><strong>Company Website:</strong> <a href={job.companyWebsite} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline break-all">{job.companyWebsite}</a></p>
              <p><strong>Job Posted Time:</strong> {new Date(job.jobPostedTime).toLocaleString()}</p>
              <p><strong>Job Link:</strong> <a href={job.jobLink} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline break-all">{job.jobLink}</a></p>
              <p><strong>Job Description:</strong></p>
              <div
                className="whitespace-pre-wrap bg-gray-100 p-3 rounded-md mt-1"
                dangerouslySetInnerHTML={{ __html: job.jobDescription }}
              />
              {renderArray(job.jobRequirements, "Requirements")}
              {renderArray(job.jobBenefits, "Benefits")}
              {renderArray(job.jobTypes, "Job Types")}
              <p><strong>Is New Job:</strong> {job.isNewJob ? "Yes" : "No"}</p>
              <p><strong>Is No Resume Job:</strong> {job.isNoResumeJob ? "Yes" : "No"}</p>
              <p><strong>Is Urgently Hiring:</strong> {job.isUrgentlyHiring ? "Yes" : "No"}</p>
              {renderOptional(job.roleType, "Role Type")}
              {renderOptional(job.minSalary, "Min Salary")}
              {renderOptional(job.maxSalary, "Max Salary")}
              {renderOptional(job.salaryGranularity, "Salary Granularity")}
              {renderOptional(job.hiresNeeded, "Hires Needed")}
              {renderOptional(job.zip, "ZIP Code")}
              {renderOptional(job.placeId, "Place ID")}
              {renderOptional(job.latitude, "Latitude")}
              {renderOptional(job.longitude, "Longitude")}
              <p><strong>Location Count:</strong> {job.locationCount}</p>
              {renderOptional(job.facebook, "Facebook")}
              {renderOptional(job.instagram, "Instagram")}
              {renderOptional(job.tiktok, "TikTok")}
              {renderOptional(job.youtube, "YouTube")}
              {renderOptional(job.twitter, "Twitter")}
              {renderOptional(job.yelp, "Yelp")}
              {/* Explicitly render scheduling link as a clickable link */}
              {job.schedulingLink && (
                  <p><strong>Scheduling Link:</strong> <a href={job.schedulingLink} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline break-all">{job.schedulingLink}</a></p>
              )}
              <hr className="my-2"/>
              <p><strong>Invocation ID:</strong> {job.invocationID}</p>
              <p><strong>Task ID:</strong> {job.taskID}</p>
              <p><strong>Date Represented:</strong> {new Date(job.dateRepresented).toLocaleString()}</p>
              <p><strong>Date Collected:</strong> {new Date(job.dateCollected).toLocaleString()}</p>
              <p><strong>Attempt ID:</strong> {job.attemptID}</p>
              <p><strong>Created At:</strong> {new Date(job.created_at).toLocaleString()}</p>
              <p><strong>Updated At:</strong> {new Date(job.updated_at).toLocaleString()}</p>
            </div>
          </DialogContent>
        </Dialog>
      )
    },
  },
] 