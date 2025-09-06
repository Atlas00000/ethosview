import React from "react";
import type { TopPEResponse } from "../../types/api";

export function PELeaders({ top }: { top: TopPEResponse }) {
  const items = top.top_performers;
  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <h2 className="text-xl font-semibold mb-3 text-gradient">Top P/E (lower is better)</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        {items.slice(0, 10).map((it) => (
          <div key={it.company_id} className="glass-card p-3">
            <div className="text-sm font-medium" style={{ color: "#0B2545" }}>{it.company_name}</div>
            <div className="text-xs" style={{ color: "#374151" }}>Rank #{it.rank}</div>
            <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{it.value.toFixed(2)}</div>
          </div>
        ))}
      </div>
    </section>
  );
}


