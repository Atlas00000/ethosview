"use client";
import React, { useMemo, useState } from "react";
import type { SectorComparison } from "../../types/api";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts";

const COLORS = ["#1E6AE1", "#2AB3A6", "#3986FF", "#1DAA8E", "#9AD36A", "#FFB547"];

export function SectorPie({ sectors }: { sectors: SectorComparison[] }) {
  const [metric, setMetric] = useState<"count" | "total_mcap">("count");
  const [desc, setDesc] = useState(true);
  const [expanded, setExpanded] = useState(false);
  const [activeIndex, setActiveIndex] = useState(0);

  const sorted = useMemo(() => {
    const arr = [...sectors];
    const valueOf = (s: SectorComparison) => (metric === "count" ? s.company_count : s.total_market_cap);
    arr.sort((a, b) => (desc ? valueOf(b) - valueOf(a) : valueOf(a) - valueOf(b)));
    return arr;
  }, [sectors, metric, desc]);

  const visible = useMemo(() => sorted.slice(0, expanded ? sorted.length : 10), [sorted, expanded]);
  const data = visible.map((s) => ({ name: s.sector, value: metric === "count" ? s.company_count : s.total_market_cap }));
  const total = data.reduce((acc, d) => acc + (typeof d.value === 'number' ? d.value : 0), 0);
  const active = visible[activeIndex] || visible[0];

  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-3">
        <h2 className="text-xl font-semibold text-gradient">Sector distribution</h2>
        <div className="flex items-center gap-2">
          <select className="border rounded px-2 py-1 hover-lift" value={metric} onChange={(e) => setMetric(e.target.value as any)}>
            <option value="count">By company count</option>
            <option value="total_mcap">By total market cap</option>
          </select>
          <button className="glass-card px-3 py-1 btn-sheen" onClick={() => setDesc(v => !v)}>{desc ? "Desc" : "Asc"}</button>
          <button className="glass-card px-3 py-1 btn-sheen" onClick={() => setExpanded(v => !v)}>{expanded ? "Show fewer" : "Show all"}</button>
        </div>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="glass-card p-3 hover-lift tilt-hover" style={{ width: "100%", height: 300 }}>
          <ResponsiveContainer>
            <PieChart>
              <Pie
                data={data}
                dataKey="value"
                nameKey="name"
                innerRadius={70}
                outerRadius={110}
                paddingAngle={3}
                isAnimationActive
                animationDuration={600}
                onMouseEnter={(_, idx) => setActiveIndex(idx)}
              >
                {data.map((_, idx) => (
                  <Cell key={idx} fill={COLORS[idx % COLORS.length]} opacity={activeIndex === idx ? 1 : 0.85} />
                ))}
              </Pie>
              <Tooltip formatter={(val: number, _name: string, p) => {
                const pct = total > 0 ? (val / total) * 100 : 0;
                return [metric === 'count' ? `${val} companies` : formatCompact(val), `${p?.payload?.name} (${pct.toFixed(1)}%)`];
              }} />
            </PieChart>
          </ResponsiveContainer>
        </div>
        <div className="glass-card p-4 animate-fade-in-up">
          <div className="text-sm font-medium" style={{ color: "#0B2545" }}>{active?.sector || '—'}</div>
          <div className="mt-2 grid grid-cols-2 gap-2">
            <div className="glass-card p-2 text-center">
              <div className="text-[10px]" style={{ color: "#374151" }}>Companies</div>
              <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{active?.company_count ?? '—'}</div>
            </div>
            <div className="glass-card p-2 text-center">
              <div className="text-[10px]" style={{ color: "#374151" }}>Total Cap</div>
              <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{active ? formatCompact(active.total_market_cap) : '—'}</div>
            </div>
            <div className="glass-card p-2 text-center">
              <div className="text-[10px]" style={{ color: "#374151" }}>Avg ESG</div>
              <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{active ? active.avg_esg_score.toFixed(2) : '—'}</div>
            </div>
            <div className="glass-card p-2 text-center">
              <div className="text-[10px]" style={{ color: "#374151" }}>Avg P/E</div>
              <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{active ? active.avg_pe_ratio.toFixed(2) : '—'}</div>
            </div>
          </div>
          <div className="mt-2 flex items-center justify-between text-xs">
            <div><span style={{ color: "#374151" }}>Best:</span> <span className="font-medium" style={{ color: "#0B2545" }}>{active?.best_esg_company || '—'}</span></div>
            <div><span style={{ color: "#374151" }}>Worst:</span> <span className="font-medium" style={{ color: "#0B2545" }}>{active?.worst_esg_company || '—'}</span></div>
          </div>
        </div>
      </div>
      <div className="mt-4 grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
        {visible.map((s, idx) => {
          const val = metric === 'count' ? s.company_count : s.total_market_cap;
          const pct = total > 0 ? (Number(val) / total) * 100 : 0;
          return (
            <div key={s.sector} className="group glass-card p-3 hover-lift tilt-hover animate-fade-in-up relative overflow-hidden" onMouseEnter={() => setActiveIndex(idx)}>
              <div className="absolute inset-0 pointer-events-none button-sheen opacity-0 group-hover:opacity-100" />
              <div className="flex items-center justify-between">
                <div className="text-sm font-medium" style={{ color: "#0B2545" }}>{s.sector}</div>
                <span className="px-2 py-0.5 rounded-full text-[10px] bg-white/60">{pct.toFixed(1)}%</span>
              </div>
              <div className="mt-2 grid grid-cols-2 gap-2">
                <div className="glass-card p-2 text-center">
                  <div className="text-[10px]" style={{ color: "#374151" }}>{metric === 'count' ? 'Companies' : 'Total Cap'}</div>
                  <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{metric === 'count' ? val : formatCompact(Number(val))}</div>
                </div>
                <div className="glass-card p-2 text-center">
                  <div className="text-[10px]" style={{ color: "#374151" }}>Avg ESG</div>
                  <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{s.avg_esg_score.toFixed(2)}</div>
                </div>
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


