import { api } from "../../services/api";
import type { AnalyticsSummaryResponse } from "../../types/api";
import { SectorBar } from "../../components/dashboard/SectorBar";
import { TopESGList } from "../../components/dashboard/TopESGList";

export const revalidate = 120;

export default async function DashboardPage() {
  const analytics: AnalyticsSummaryResponse = await api.analyticsSummary();
  return (
    <main className="max-w-6xl mx-auto px-4 py-10">
      <h1 className="text-2xl md:text-3xl font-semibold text-gradient">Dashboard</h1>
      <div className="mt-6 grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="glass-card p-4 lg:col-span-2">
          <h2 className="text-sm font-medium" style={{ color: "#374151" }}>Sector ESG averages</h2>
          <div className="mt-3">
            <SectorBar sectors={analytics.sector_comparisons.slice(0, 12)} />
          </div>
        </div>
        <div className="glass-card p-4">
          <h2 className="text-sm font-medium" style={{ color: "#374151" }}>Top ESG performers</h2>
          <div className="mt-3">
            <TopESGList items={analytics.top_esg_performers.slice(0, 10)} />
          </div>
        </div>
      </div>
    </main>
  );
}
