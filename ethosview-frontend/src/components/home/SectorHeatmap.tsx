"use client";
import React, { useMemo, useState } from "react";
import type { SectorComparison } from "../../types/api";

export function SectorHeatmap({ sectors }: { sectors: SectorComparison[] }) {
  if (!sectors.length) return null;

  const [metric, setMetric] = useState<"esg" | "pe" | "mcap" | "total_mcap">("esg");
  const [desc, setDesc] = useState(true);

  const sorted = useMemo(() => {
    const valueOf = (s: SectorComparison) => {
      if (metric === "esg") return s.avg_esg_score;
      if (metric === "pe") return s.avg_pe_ratio;
      if (metric === "mcap") return s.avg_market_cap;
      return s.total_market_cap;
    };
    const arr = [...sectors];
    arr.sort((a, b) => (desc ? valueOf(b) - valueOf(a) : valueOf(a) - valueOf(b)));
    return arr;
  }, [sectors, metric, desc]);

  const values = sorted.map((s) => metric === "esg" ? s.avg_esg_score : metric === "pe" ? s.avg_pe_ratio : metric === "mcap" ? s.avg_market_cap : s.total_market_cap);
  const min = Math.min(...values);
  const max = Math.max(...values);
  const norm = (v: number) => (max === min ? 0.5 : (v - min) / (max - min));

  return (
    <section className="max-w-6xl mx-auto px-4 py-12 relative overflow-hidden">
      <div className="blob" style={{ top: -70, right: -50, width: 220, height: 220, background: "rgba(42,179,166,0.12)", borderRadius: 9999 }} />
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold text-gradient">Sector ESG heatmap</h2>
        <div className="flex items-center gap-2">
          <select className="border rounded px-2 py-1 hover-lift" value={metric} onChange={(e) => setMetric(e.target.value as any)}>
            <option value="esg">Avg ESG</option>
            <option value="pe">Avg P/E</option>
            <option value="mcap">Avg Market Cap</option>
            <option value="total_mcap">Total Market Cap</option>
          </select>
          <button className="glass-card px-3 py-1 btn-sheen" onClick={() => setDesc((v) => !v)}>{desc ? "Desc" : "Asc"}</button>
        </div>
      </div>
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 animate-fade-in-up">
        {sorted.map((s) => {
          const value = metric === "esg" ? s.avg_esg_score : metric === "pe" ? s.avg_pe_ratio : metric === "mcap" ? s.avg_market_cap : s.total_market_cap;
          const t = norm(value);
          const bg = `rgba(16, 185, 129, ${Math.max(0.15, t)})`;
          return (
            <div key={s.sector} className="group glass-card p-3 hover-lift tilt-hover relative overflow-hidden" style={{ backgroundColor: bg }} title={`${s.sector}`}>
              <div className="absolute inset-0 pointer-events-none button-sheen opacity-0 group-hover:opacity-100" />
              <div className="flex items-center justify-between">
                <div className="text-sm font-medium" style={{ color: "#0B2545" }}>{s.sector}</div>
                <span className="px-2 py-0.5 rounded-full text-[10px] bg-white/50">{s.company_count} cos</span>
              </div>
              <div className="mt-2 grid grid-cols-3 gap-2">
                <div className="glass-card p-2 text-center">
                  <div className="text-[10px]" style={{ color: "#374151" }}>Avg ESG</div>
                  <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{s.avg_esg_score.toFixed(2)}</div>
                </div>
                <div className="glass-card p-2 text-center">
                  <div className="text-[10px]" style={{ color: "#374151" }}>Avg P/E</div>
                  <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{s.avg_pe_ratio.toFixed(2)}</div>
                </div>
                <div className="glass-card p-2 text-center">
                  <div className="text-[10px]" style={{ color: "#374151" }}>Avg Cap</div>
                  <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{formatCompact(s.avg_market_cap)}</div>
                </div>
              </div>
              <div className="mt-2 h-1.5 bg-white/30 rounded">
                <div className="h-1.5 rounded" style={{ width: `${Math.round(t * 100)}%`, background: '#1E6AE1', transition: 'width 320ms var(--ease-enter)' }} />
              </div>
              <div className="mt-2 flex items-center justify-between text-xs">
                <div><span style={{ color: "#374151" }}>Best:</span> <span className="font-medium" style={{ color: "#0B2545" }}>{s.best_esg_company || '—'}</span></div>
                <div><span style={{ color: "#374151" }}>Worst:</span> <span className="font-medium" style={{ color: "#0B2545" }}>{s.worst_esg_company || '—'}</span></div>
              </div>
            </div>
          );
        })}
      </div>
    </section>
  );
}

function formatCompact(n: number) {
  if (!isFinite(n)) return '—';
  if (n >= 1e12) return (n / 1e12).toFixed(2) + 'T';
  if (n >= 1e9) return (n / 1e9).toFixed(2) + 'B';
  if (n >= 1e6) return (n / 1e6).toFixed(2) + 'M';
  if (n >= 1e3) return (n / 1e3).toFixed(2) + 'K';
  return n.toFixed(0);
}


