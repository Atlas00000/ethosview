"use client";
import React from "react";
import type { PerformanceMetric } from "../../types/api";

export function ESGHighlights({ top }: { top: PerformanceMetric[] }) {
  const items = top.slice(0, 5);

  return (
    <section className="max-w-6xl mx-auto px-4 py-12 relative overflow-hidden">
      <div className="blob" style={{ top: -60, left: -50, width: 200, height: 200, background: "rgba(30,106,225,0.1)", borderRadius: 9999 }} />
      <h2 className="text-xl font-semibold mb-3 text-gradient">ESG highlights</h2>
      <div className="glass-card divide-y animate-fade-in-up">
        {items.map((it) => (
          <div key={it.company_id} className="p-3 flex items-center justify-between hover-lift">
            <div className="truncate pr-3">
              <div className="text-sm font-medium truncate" style={{ color: "#0B2545" }}>{it.company_name}</div>
              <div className="text-xs" style={{ color: "#374151" }}>Rank #{it.rank} • Percentile {it.percentile.toFixed(1)}%</div>
            </div>
            <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{it.value.toFixed(2)}</div>
          </div>
        ))}
      </div>
    </section>
  );
}


