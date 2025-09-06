import Link from "next/link";

export default function NotFound() {
  return (
    <div className="min-h-[60vh] flex items-center justify-center">
      <div className="glass-card p-8 text-center shadow-elevated">
        <h1 className="text-3xl font-semibold text-gradient">404</h1>
        <p className="mt-2" style={{ color: "#374151" }}>
          The page you’re looking for doesn’t exist.
        </p>
        <Link href="#hero" className="inline-block mt-4 btn-primary rounded px-4 py-2">
          Go home
        </Link>
      </div>
    </div>
  );
}


