"use client";
import React from "react";
import type { CorrelationResponse } from "../../types/api";
import { ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip, Cell } from "recharts";
import { useCountUp } from "./useCountUp";

function formatValue(v: unknown): string {
  const num = typeof v === "number" && Number.isFinite(v) ? v : NaN;
  return Number.isFinite(num) ? num.toFixed(2) : "—";
}

export function CorrelationTeaser({ corr }: { corr: CorrelationResponse | Record<string, any> }) {
  const metrics = [
    { key: "esg_market_cap_corr", label: "ESG vs Market Cap" },
    { key: "esg_pe_corr", label: "ESG vs P/E" },
    { key: "esg_roe_corr", label: "ESG vs ROE" },
    { key: "esg_profit_corr", label: "ESG vs Profit" },
  ] as const;

  const bars = metrics.map((m) => ({ name: m.label.split(" ")[2], value: Math.abs((corr as any)?.[m.key] ?? 0), raw: (corr as any)?.[m.key] }));

  const sample = (corr as any)?.sample_size ?? 0;
  const avgESG = (corr as any)?.avg_esg_score ?? null;
  const avgPE = (corr as any)?.avg_pe_ratio ?? null;
  const avgROE = (corr as any)?.avg_roe ?? null;
  const avgProfit = (corr as any)?.avg_profit_margin ?? null;

  const cMarket = useCountUp(typeof (corr as any)?.esg_market_cap_corr === 'number' ? (corr as any).esg_market_cap_corr : 0);
  const cPE = useCountUp(typeof (corr as any)?.esg_pe_corr === 'number' ? (corr as any).esg_pe_corr : 0);
  const cROE = useCountUp(typeof (corr as any)?.esg_roe_corr === 'number' ? (corr as any).esg_roe_corr : 0);
  const cProfit = useCountUp(typeof (corr as any)?.esg_profit_corr === 'number' ? (corr as any).esg_profit_corr : 0);
  const avgMktCap = (corr as any)?.avg_market_cap ?? null;

  const signColor = (v: number | null | undefined) => (typeof v === 'number' ? (v > 0 ? '#1E6AE1' : v < 0 ? '#E11D48' : '#0B2545') : '#0B2545');
  const strength = (r: number | null | undefined) => {
    if (typeof r !== 'number' || !isFinite(r)) return '—';
    const a = Math.abs(r);
    if (a >= 0.7) return 'Strong';
    if (a >= 0.4) return 'Moderate';
    if (a > 0) return 'Weak';
    return '—';
  };

  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <h2 className="text-xl font-semibold mb-3 text-gradient">ESG ↔ Financial correlation</h2>

      <div className="grid grid-cols-2 sm:grid-cols-5 gap-3 mb-3">
        <div className="glass-card p-3 text-center hover-lift">
          <div className="text-[11px]" style={{ color: "#374151" }}>Sample size</div>
          <div className="text-lg font-semibold" style={{ color: "#0B2545" }}>{Number.isFinite(sample) ? sample : '—'}</div>
        </div>
        <div className="glass-card p-3 text-center hover-lift">
          <div className="text-[11px]" style={{ color: "#374151" }}>Avg ESG</div>
          <div className="text-lg font-semibold" style={{ color: "#1D9A6C" }}>{formatValue(avgESG)}</div>
        </div>
        <div className="glass-card p-3 text-center hover-lift">
          <div className="text-[11px]" style={{ color: "#374151" }}>Avg P/E</div>
          <div className="text-lg font-semibold" style={{ color: "#0B2545" }}>{formatValue(avgPE)}</div>
        </div>
        <div className="glass-card p-3 text-center hover-lift">
          <div className="text-[11px]" style={{ color: "#374151" }}>Avg ROE</div>
          <div className="text-lg font-semibold" style={{ color: "#0B2545" }}>{formatValue(avgROE)}</div>
        </div>
        <div className="glass-card p-3 text-center hover-lift">
          <div className="text-[11px]" style={{ color: "#374151" }}>Avg MktCap</div>
          <div className="text-lg font-semibold" style={{ color: "#0B2545" }}>{formatValue(avgMktCap)}</div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div className="glass-card p-3 animate-fade-in-up">
          <div className="text-xs mb-1" style={{ color: "#374151" }}>Correlation magnitude (|r|)</div>
          <div style={{ width: "100%", height: 160 }}>
            <ResponsiveContainer>
              <BarChart data={bars} margin={{ top: 6, right: 8, left: 0, bottom: 0 }}>
                <XAxis dataKey="name" tick={{ fontSize: 10 }} />
                <YAxis domain={[0, 1]} tick={{ fontSize: 10 }} />
                <Tooltip formatter={(v: any, n: any, p: any) => [formatValue(v), `${p.payload?.name}`]} />
                <Bar dataKey="value" radius={[6, 6, 0, 0]} className="hover-lift">
                  {bars.map((b, idx) => (
                    <Cell key={`c-${idx}`} fill={signColor(b.raw)} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
        <div className="glass-card p-3 animate-fade-in-up">
          <div className="text-xs mb-1" style={{ color: "#374151" }}>Direction and value</div>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-2">
            <div className="glass-card p-2 text-center tilt-hover">
              <div className="text-[10px]" style={{ color: "#374151" }}>Mkt Cap</div>
              <div className="text-base font-semibold" style={{ color: signColor((corr as any)?.esg_market_cap_corr) }}>{formatValue(cMarket)} {((corr as any)?.esg_market_cap_corr ?? 0) > 0 ? '↑' : ((corr as any)?.esg_market_cap_corr ?? 0) < 0 ? '↓' : ''}</div>
            </div>
            <div className="glass-card p-2 text-center tilt-hover">
              <div className="text-[10px]" style={{ color: "#374151" }}>P/E</div>
              <div className="text-base font-semibold" style={{ color: signColor((corr as any)?.esg_pe_corr) }}>{formatValue(cPE)} {((corr as any)?.esg_pe_corr ?? 0) > 0 ? '↑' : ((corr as any)?.esg_pe_corr ?? 0) < 0 ? '↓' : ''}</div>
            </div>
            <div className="glass-card p-2 text-center tilt-hover">
              <div className="text-[10px]" style={{ color: "#374151" }}>ROE</div>
              <div className="text-base font-semibold" style={{ color: signColor((corr as any)?.esg_roe_corr) }}>{formatValue(cROE)} {((corr as any)?.esg_roe_corr ?? 0) > 0 ? '↑' : ((corr as any)?.esg_roe_corr ?? 0) < 0 ? '↓' : ''}</div>
            </div>
            <div className="glass-card p-2 text-center tilt-hover">
              <div className="text-[10px]" style={{ color: "#374151" }}>Profit</div>
              <div className="text-base font-semibold" style={{ color: signColor((corr as any)?.esg_profit_corr) }}>{formatValue(cProfit)} {((corr as any)?.esg_profit_corr ?? 0) > 0 ? '↑' : ((corr as any)?.esg_profit_corr ?? 0) < 0 ? '↓' : ''}</div>
            </div>
          </div>
          <div className="mt-2 flex flex-wrap gap-2">
            {metrics.map((m) => (
              <span key={m.key} className="px-2 py-0.5 rounded-full text-xs bg-white/40">
                {m.label}: {strength((corr as any)?.[m.key])}
              </span>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}


