"use client";
import React, { useMemo, useState } from "react";
import type { AdvancedSummaryResponse } from "../../types/api";
import { useCountUp } from "./useCountUp";

export function AdvancedInsights({ summary }: { summary: AdvancedSummaryResponse }) {
  const s = summary?.summary;
  const [expanded, setExpanded] = useState(false);

  const assessed = useCountUp(s?.risk_summary?.total_companies_assessed ?? 0);
  const avgRisk = useCountUp((s?.risk_summary?.average_risk_score as number) ?? 0);

  const riskDist = s?.risk_summary?.risk_distribution || { low: 0, medium: 0, high: 0 };
  const riskTotal = (riskDist.low ?? 0) + (riskDist.medium ?? 0) + (riskDist.high ?? 0);

  const esgTr = s?.trend_summary?.esg_trends || { improving: 0, declining: 0, stable: 0 };
  const priceTr = s?.trend_summary?.price_trends || { up: 0, down: 0, stable: 0 };
  const esgTotal = (esgTr.improving ?? 0) + (esgTr.declining ?? 0) + (esgTr.stable ?? 0);
  const priceTotal = (priceTr.up ?? 0) + (priceTr.down ?? 0) + (priceTr.stable ?? 0);

  const optReady = Boolean(s?.portfolio_optimization);

  return (
    <section className="max-w-6xl mx-auto px-4 py-8 relative overflow-hidden">
      <div className="blob" style={{ bottom: -70, right: -50, width: 200, height: 200, background: "rgba(42,179,166,0.10)", borderRadius: 9999 }} />
      <div className="flex items-center justify-between mb-3">
        <h2 className="text-xl font-semibold text-gradient">Advanced insights</h2>
        <button className="glass-card px-3 py-1 btn-sheen" onClick={() => setExpanded(v => !v)}>{expanded ? "Show fewer" : "Show more"}</button>
      </div>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-3">
        <div className="glass-card p-4 hover-lift tilt-hover animate-fade-in-up">
          <div className="text-xs" style={{ color: "#374151" }}>Assessed companies</div>
          <div className="text-2xl font-semibold" style={{ color: "#0B2545" }}>{Math.round(assessed).toLocaleString()}</div>
        </div>
        <div className="glass-card p-4 hover-lift tilt-hover animate-fade-in-up">
          <div className="text-xs" style={{ color: "#374151" }}>Average risk score</div>
          <div className="text-2xl font-semibold" style={{ color: "#0B2545" }}>{avgRisk.toFixed(2)}</div>
          <div className="mt-2 h-1.5 bg-white/30 rounded">
            <div className="h-1.5 rounded" style={{ width: `${Math.max(0, Math.min(100, (Number(s?.risk_summary?.average_risk_score) || 0) * 10))}%`, background: '#FFB547', transition: 'width 320ms var(--ease-enter)' }} />
          </div>
        </div>
        <div className="glass-card p-4 hover-lift tilt-hover animate-fade-in-up">
          <div className="text-xs" style={{ color: "#374151" }}>Portfolio optimization</div>
          <div className="text-sm font-semibold mt-1" style={{ color: optReady ? '#1D9A6C' : '#374151' }}>{optReady ? 'Ready' : '—'}</div>
          <div className="mt-2 flex items-center gap-2">
            <span className={`px-2 py-0.5 rounded-full text-[10px] ${optReady ? 'bg-white/60' : 'bg-white/40'}`}>{optReady ? 'Optimization available' : 'Awaiting data'}</span>
          </div>
        </div>
      </div>

      <div className="mt-4 grid grid-cols-1 lg:grid-cols-2 gap-3">
        {/* Risk distribution */}
        <div className="glass-card p-4 animate-fade-in-up">
          <div className="text-sm font-medium mb-2" style={{ color: "#0B2545" }}>Risk distribution</div>
          <div className="grid grid-cols-3 gap-2">
            {(['low','medium','high'] as const).map((k) => {
              const v = (riskDist as any)[k] ?? 0;
              const pct = riskTotal > 0 ? (v / riskTotal) * 100 : 0;
              const color = k === 'low' ? '#1D9A6C' : k === 'medium' ? '#FFB547' : '#E25555';
              return (
                <div key={k} className="glass-card p-2 text-center hover-lift">
                  <div className="text-[10px]" style={{ color: '#374151' }}>{k.toUpperCase()}</div>
                  <div className="text-sm font-semibold" style={{ color: '#0B2545' }}>{v}</div>
                  <div className="mt-1 h-1.5 bg-white/30 rounded">
                    <div className="h-1.5 rounded" style={{ width: `${Math.round(pct)}%`, background: color, transition: 'width 320ms var(--ease-enter)' }} />
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Trend summary */}
        <div className="glass-card p-4 animate-fade-in-up">
          <div className="text-sm font-medium mb-2" style={{ color: "#0B2545" }}>Trend summary</div>
          <div className="grid grid-cols-2 gap-2">
            <div className="glass-card p-2">
              <div className="text-xs mb-1" style={{ color: '#374151' }}>ESG trends</div>
              <TrendRow label="Improving" value={esgTr.improving ?? 0} total={esgTotal} color="#1D9A6C" />
              <TrendRow label="Stable" value={esgTr.stable ?? 0} total={esgTotal} color="#3986FF" />
              <TrendRow label="Declining" value={esgTr.declining ?? 0} total={esgTotal} color="#E25555" />
            </div>
            <div className="glass-card p-2">
              <div className="text-xs mb-1" style={{ color: '#374151' }}>Price trends</div>
              <TrendRow label="Up" value={priceTr.up ?? 0} total={priceTotal} color="#1D9A6C" />
              <TrendRow label="Stable" value={priceTr.stable ?? 0} total={priceTotal} color="#3986FF" />
              <TrendRow label="Down" value={priceTr.down ?? 0} total={priceTotal} color="#E25555" />
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

function TrendRow({ label, value, total, color }: { label: string; value: number; total: number; color: string }) {
  const pct = total > 0 ? (value / total) * 100 : 0;
  return (
    <div className="mt-1">
      <div className="flex items-center justify-between text-xs" style={{ color: '#374151' }}>
        <span>{label}</span>
        <span>{value}</span>
      </div>
      <div className="mt-1 h-1.5 bg-white/30 rounded">
        <div className="h-1.5 rounded" style={{ width: `${Math.round(pct)}%`, background: color, transition: 'width 300ms var(--ease-enter)' }} />
      </div>
    </div>
  );
}


