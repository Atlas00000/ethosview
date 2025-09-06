"use client";
import React from "react";
import type { SectorComparison } from "../../types/api";

export function SectorHeatmap({ sectors }: { sectors: SectorComparison[] }) {
  if (!sectors.length) return null;

  const scores = sectors.map((s) => s.avg_esg_score);
  const min = Math.min(...scores);
  const max = Math.max(...scores);
  const norm = (v: number) => (max === min ? 0.5 : (v - min) / (max - min));

  return (
    <section className="max-w-6xl mx-auto px-4 py-12 relative overflow-hidden">
      <div className="blob" style={{ top: -70, right: -50, width: 220, height: 220, background: "rgba(42,179,166,0.12)", borderRadius: 9999 }} />
      <h2 className="text-xl font-semibold mb-4 text-gradient">Sector ESG heatmap</h2>
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 animate-fade-in-up">
        {sectors.map((s) => {
          const t = norm(s.avg_esg_score);
          const bg = `rgba(16, 185, 129, ${Math.max(0.15, t)})`;
          return (
            <div key={s.sector} className="glass-card p-3 hover-lift" style={{ backgroundColor: bg }}>
              <div className="text-sm font-medium" style={{ color: "#0B2545" }}>{s.sector}</div>
              <div className="text-xs" style={{ color: "#0B2545" }}>Avg ESG: {s.avg_esg_score.toFixed(2)}</div>
              <div className="text-xs" style={{ color: "#374151" }}>Companies: {s.company_count}</div>
            </div>
          );
        })}
      </div>
    </section>
  );
}


