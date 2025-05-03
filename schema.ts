/**
 * Key fields described below:
 * - invocationID can be used to distinguish different ingestion attempts,
 * - dateRepresented can be used to represent the date of the data point/segment by date,
 */
export interface Row extends RowDefinition {
  companyName: string; // Company name
  jobPostedTime: Date; // Job posting creation date
  companyRating: number; // Company rating
  jobTitle: string; // Job title
  roleType?: string | null; // Role type
  minSalary?: number; // Minimum salary
  maxSalary?: number; // Maximum salary
  salaryGranularity?: string; // Salary granularity
  companyAddress: string; // Full address
  city: string; // City of the job
  state?: string; // State of the job
  zip?: string; // Zip code
  jobID: string; // Job ID
  isNewJob: boolean; // Is this a new job
  jobLink: string; // Job link
  hiresNeeded?: string; // Number of positions to fill
  isNoResumeJob: boolean; // Whether job requires resume
  jobRequirements: string[]; // Job requirements
  locationCount: number; // Number of locations
  jobDescription: string; // Job description snippet
  jobBenefits: string[]; // Job benefits
  isUrgentlyHiring: boolean; // Urgently hiring flag
  jobTypes: string[]; // Job types
  companyWebsite: string; // Company website link
  facebook?: string; // Facebook link
  instagram?: string; // Instagram link
  tiktok?: string; // TikTok link
  youtube?: string; // YouTube link
  twitter?: string; // Twitter link
  yelp?: string; // Yelp link
  schedulingLink?: string | null; // Scheduling/booking link
  placeId?: string; // Google Place ID
  latitude?: number; // Location latitude
  longitude?: number; // Location longitude

  invocationID: string; // Unique identifier for the delivery
  taskID: string; // Unique identifier for the task
  dateRepresented: Date; // Date of the data point
  dateCollected: Date; // Data the data point was collected
  attemptID: string; // Unique identifier for the attempt
}
