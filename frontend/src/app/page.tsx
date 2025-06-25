import Link from 'next/link';

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center gap-6 p-4">
      <h1 className="text-3xl font-bold">Welcome to System Portal</h1>
      <p className="text-center text-muted-foreground max-w-md">
        This demo shows a simplified frontend for managing the System Portal API.
        You can browse around without logging in, or sign in to access protected features.
      </p>
      <Link
        href="/login"
        className="px-6 py-3 rounded-md bg-primary text-primary-foreground hover:bg-primary/90"
      >
        Go to Login
      </Link>
    </main>
  );
}
