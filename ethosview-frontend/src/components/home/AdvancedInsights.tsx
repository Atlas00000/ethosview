import React from "react";
import type { AdvancedSummaryResponse } from "../../types/api";

export function AdvancedInsights({ summary }: { summary: AdvancedSummaryResponse }) {
  const s = summary?.summary;
  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <h2 className="text-xl font-semibold mb-3 text-gradient">Advanced insights</h2>
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div className="glass-card p-3">
          <div className="text-xs" style={{ color: "#374151" }}>Assessed</div>
          <div className="text-lg font-medium" style={{ color: "#0B2545" }}>{s?.risk_summary?.total_companies_assessed ?? 0}</div>
        </div>
        <div className="glass-card p-3">
          <div className="text-xs" style={{ color: "#374151" }}>Avg risk</div>
          <div className="text-lg font-medium" style={{ color: "#0B2545" }}>{(s?.average_risk_score as any) ?? 0}</div>
        </div>
        <div className="glass-card p-3">
          <div className="text-xs" style={{ color: "#374151" }}>Optimization</div>
          <div className="text-lg font-medium" style={{ color: "#0B2545" }}>{s?.portfolio_optimization ? "Available" : "—"}</div>
        </div>
      </div>
    </section>
  );
}


