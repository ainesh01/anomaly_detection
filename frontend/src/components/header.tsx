import Link from 'next/link';
import { Button } from '@/components/ui/button';

export function Header() {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 max-w-screen-2xl items-center justify-center">
        <nav className="flex items-center space-x-4 lg:space-x-6">
          <Link href="/" passHref>
            <Button variant="ghost">Home</Button>
          </Link>
          <Link href="/data" passHref>
            <Button variant="ghost">Data</Button>
          </Link>
          <Link href="/anomalies" passHref>
            <Button variant="ghost">Anomalies</Button>
          </Link>
          <Link href="/rules" passHref>
            <Button variant="ghost">Rules</Button>
          </Link>
        </nav>
      </div>
    </header>
  );
} 