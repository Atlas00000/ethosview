import React from "react";
import type { ESGScoresListResponse } from "../../types/api";

export function ESGFeed({ list }: { list: ESGScoresListResponse }) {
  const items = list?.scores || [];
  if (!items.length) return null;
  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <h2 className="text-xl font-semibold mb-3 text-gradient">Latest ESG scores</h2>
      <div className="glass-card divide-y">
        {items.map((s) => (
          <div key={s.id} className="p-3 flex items-center justify-between">
            <div className="truncate pr-3">
              <div className="text-sm font-medium" style={{ color: "#0B2545" }}>{s.company_name || s.company_id}</div>
              <div className="text-xs" style={{ color: "#374151" }}>{new Date(s.score_date).toISOString().slice(0,10)}</div>
            </div>
            <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{s.overall_score.toFixed(2)}</div>
          </div>
        ))}
      </div>
    </section>
  );
}


