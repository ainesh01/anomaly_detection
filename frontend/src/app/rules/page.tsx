"use client"; // Make this a client component

import * as React from "react"; // Import React for useState
import { useState, useCallback } from "react"; // Import useState and useCallback
import { DataTable } from "@/components/ui/data-table";
import { Button } from "@/components/ui/button";
import type { ColumnDef } from "@tanstack/react-table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"; // Import Dialog components
import { Input } from "@/components/ui/input"; // Import Input
import { Label } from "@/components/ui/label"; // Import Label
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"; // Import Select components
import { Switch } from "@/components/ui/switch"; // Import Switch
import { toast } from "sonner"; // Import toast from sonner

// Define a basic type for the Rule data (adjust based on actual API response)
type Rule = {
  id: string; // Assuming ID is returned, adjust if needed
  name: string;
  type: string; // Changed from metric
  operator: string; // Changed from comparison_operator
  value: number; // Changed from threshold
  is_active: boolean; // Add isActive field
};

export default function RulesPage() { // Changed to non-async as fetching is moved
  const [rulesData, setRulesData] = useState<Rule[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  // State for the form fields
  const [ruleName, setRuleName] = useState("");
  const [ruleDescription, setRuleDescription] = useState("");
  const [ruleMetric, setRuleMetric] = useState<string | undefined>(undefined);
  const [ruleOperator, setRuleOperator] = useState<string | undefined>(undefined);
  const [ruleThreshold, setRuleThreshold] = useState<string>(""); // Keep as string for input control

  // Add a state for submission feedback
  const [submissionStatus, setSubmissionStatus] = useState<'idle' | 'submitting' | 'success' | 'error'>('idle');
  const [submissionError, setSubmissionError] = useState<string | null>(null);

  // Define fetchRules outside useEffect using useCallback
  const fetchRules = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch('http://localhost:8080/api/anomaly-rules');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      // Ensure ID is a string and isActive is boolean
      const formattedData = data.map((rule: any) => ({
        ...rule,
        id: String(rule.id), // Convert ID to string if necessary
        is_active: !!rule.is_active, // Ensure boolean type
      }));
      setRulesData(formattedData);
    } catch (e) {
      if (e instanceof Error) {
        setError(`Failed to fetch rules: ${e.message}`);
      } else {
        setError("An unknown error occurred");
      }
      console.error("Fetch error:", e);
    } finally {
      setIsLoading(false);
    }
  }, []); // Empty dependency array for useCallback as it doesn't depend on props/state

  React.useEffect(() => {
    fetchRules();
  }, [fetchRules]); // Add fetchRules to dependency array

  // Function to handle toggling the rule's active state
  const handleToggleRule = async (ruleId: string, currentStatus: boolean) => {
    const newStatus = !currentStatus;

    // Optimistically update the UI
    setRulesData((prevRules) =>
      prevRules.map((rule) =>
        rule.id === ruleId ? { ...rule, is_active: newStatus } : rule
      )
    );

    try {
      const response = await fetch(`http://localhost:8080/api/anomaly-rules/${ruleId}/toggle`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ is_active: newStatus }),
      });

      if (!response.ok) {
        throw new Error(`Failed to toggle rule status. Status: ${response.status}`);
      }

      toast.success("Rule status updated."); // Use sonner toast

    } catch (error) {
      console.error("Failed to toggle rule:", error);
      toast.error("Error updating rule status", { // Use sonner toast
        description: error instanceof Error ? error.message : "Unknown error",
      });

      // Revert the optimistic update on error
      setRulesData((prevRules) =>
        prevRules.map((rule) =>
          rule.id === ruleId ? { ...rule, is_active: currentStatus } : rule
        )
      );
    }
  };

  // Define columns for the DataTable including the toggle switch
  const columns: ColumnDef<Rule>[] = [
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "type",
      header: "Metric",
    },
    {
      accessorKey: "operator",
      header: "Operator",
    },
    {
      accessorKey: "value",
      header: "Threshold",
    },
    {
      accessorKey: "is_active",
      header: "Active",
      cell: ({ row }) => {
        const rule = row.original;
        return (
          <Switch
            checked={rule.is_active}
            onCheckedChange={() => handleToggleRule(rule.id, rule.is_active)}
            aria-label={rule.is_active ? "Deactivate rule" : "Activate rule"}
          />
        );
      },
    },
  ];

  // Function to handle saving the rule
  const handleSaveRule = async () => {
    // Basic validation
    if (!ruleName || !ruleDescription || !ruleMetric || !ruleOperator || !ruleThreshold) {
      setSubmissionError("All fields are required.");
      setSubmissionStatus('error');
      return;
    }

    const thresholdValue = parseFloat(ruleThreshold);
    if (isNaN(thresholdValue)) {
      setSubmissionError("Threshold must be a valid number.");
      setSubmissionStatus('error');
      return;
    }

    setSubmissionStatus('submitting');
    setSubmissionError(null);

    const rulePayload = {
      Name: ruleName,
      Description: ruleDescription,
      Type: ruleMetric, // Backend expects 'Type'
      Operator: ruleOperator,
      Value: thresholdValue,
      IsActive: false, // Default to inactive for new rules
    };

    try {
      const response = await fetch('http://localhost:8080/api/anomaly-rules', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(rulePayload),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: 'Failed to parse error response' }));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      // Success
      setSubmissionStatus('success');
      // Reset form and close modal after a short delay
      setTimeout(() => {
        setRuleName("");
        setRuleDescription("");
        setRuleMetric(undefined);
        setRuleOperator(undefined);
        setRuleThreshold("");
        setIsModalOpen(false); // Close the modal
        setSubmissionStatus('idle'); // Reset status
        fetchRules(); // Refresh the rules list
      }, 1500); // Show success message for 1.5 seconds

    } catch (e) {
      let errorMessage = "An unknown error occurred during submission.";
      if (e instanceof Error) {
        errorMessage = e.message;
      }
      console.error("Submission error:", e);
      setSubmissionError(errorMessage);
      setSubmissionStatus('error');
    }
  };

  if (isLoading) {
    return <div className="container mx-auto py-10 text-center">Loading rules...</div>;
  }

  if (error) {
    return <div className="container mx-auto py-10 text-center text-red-500">Error: {error}</div>;
  }

  return (
    <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
      <div className="container mx-auto py-10">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-2xl font-bold">Anomaly Detection Rules</h1>
          <DialogTrigger asChild>
            <Button>Add New Rule</Button>
          </DialogTrigger>
        </div>
        <DataTable columns={columns} data={rulesData} />
      </div>

      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Add New Anomaly Rule</DialogTitle>
          <DialogDescription>
            Define the parameters for the new rule here. Click save when you're done.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          {/* Basic form fields - replace with actual form components later */}
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="name" className="text-right">Name</Label>
            <Input
              id="name"
              placeholder="e.g., High Max Salary"
              className="col-span-3"
              value={ruleName}
              onChange={(e) => setRuleName(e.target.value)}
            />
          </div>
          {/* Add Description Field */}
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="description" className="text-right">Description</Label>
            <Input
              id="description"
              placeholder="e.g., Alert when max salary exceeds 150k"
              className="col-span-3"
              value={ruleDescription}
              onChange={(e) => setRuleDescription(e.target.value)}
            />
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="metric" className="text-right">Metric</Label>
            <Select value={ruleMetric} onValueChange={setRuleMetric}>
              <SelectTrigger id="metric" className="col-span-3">
                <SelectValue placeholder="Select metric" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="max_salary">Max Salary</SelectItem>
                <SelectItem value="min_salary">Min Salary</SelectItem>
                <SelectItem value="company_rating">Company Rating</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="comparison_operator" className="text-right">Operator</Label>
            <Select value={ruleOperator} onValueChange={setRuleOperator}>
              <SelectTrigger id="comparison_operator" className="col-span-3">
                <SelectValue placeholder="Select operator" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value=">">Greater than (&gt;)</SelectItem>
                <SelectItem value="<">Less than (&lt;)</SelectItem>
                <SelectItem value="=">Equal to (=)</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="threshold" className="text-right">Threshold</Label>
            <Input
              id="threshold"
              type="number"
              placeholder="e.g., 150000"
              className="col-span-3"
              value={ruleThreshold}
              onChange={(e) => setRuleThreshold(e.target.value)}
            />
          </div>
        </div>
        {/* Display submission status */}
        {submissionStatus === 'error' && (
          <p className="text-red-500 text-sm mt-2">Error: {submissionError || 'Failed to save rule.'}</p>
        )}
        {submissionStatus === 'success' && (
          <p className="text-green-500 text-sm mt-2">Rule saved successfully!</p>
        )}
        <DialogFooter>
          <Button
            type="submit"
            disabled={submissionStatus === 'submitting'}
            onClick={handleSaveRule}
          >
            {submissionStatus === 'submitting' ? 'Saving...' : 'Save Rule'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
} 