import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";

export default function LandingPage() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen py-2">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl font-bold text-center">Anomaly Detection Dashboard</CardTitle>
          <CardDescription className="text-center">
            Navigate to view data, detected anomalies, or manage rules.
          </CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col space-y-4">
          <Link href="/data" passHref>
            <Button className="w-full">View Data</Button>
          </Link>
          <Link href="/anomalies" passHref>
            <Button className="w-full">View Anomalies</Button>
          </Link>
          <Link href="/rules" passHref>
            <Button className="w-full">Manage Rules</Button>
          </Link>
        </CardContent>
      </Card>
    </div>
  );
}
