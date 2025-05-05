export type JobData = {
  // Company Information
  companyName: string;
  companyRating: number;
  companyAddress: string;
  companyWebsite: string;

  // Job Information
  jobTitle: string;
  jobPostedTime: string; // Assuming CustomTime/time.Time serializes to ISO string
  jobID: string;
  jobLink: string;
  jobDescription: string;
  jobRequirements: string[];
  jobBenefits: string[];
  jobTypes: string[];
  isNewJob: boolean;
  isNoResumeJob: boolean;
  isUrgentlyHiring: boolean;

  // Role Information
  roleType?: string | null;
  minSalary?: number | null;
  maxSalary?: number | null;
  salaryGranularity?: string | null;
  hiresNeeded?: string | null;

  // Location Information
  city: string;
  state?: string | null;
  zip?: string | null;
  placeId?: string | null; // Matches json:"placeId" tag
  latitude?: number | null;
  longitude?: number | null;
  locationCount: number;

  // Social Media Links
  facebook?: string | null;
  instagram?: string | null;
  tiktok?: string | null;
  youtube?: string | null;
  twitter?: string | null;
  yelp?: string | null;

  // Additional Information
  schedulingLink?: string | null;

  // Metadata
  invocationID: string;
  taskID: string;
  dateRepresented: string; // Assuming CustomTime/time.Time serializes to ISO string
  dateCollected: string; // Assuming CustomTime/time.Time serializes to ISO string
  attemptID: string;

  // Database timestamps
  created_at: string; // Assuming CustomTime/time.Time serializes to ISO string
  updated_at: string; // Assuming CustomTime/time.Time serializes to ISO string
}; 