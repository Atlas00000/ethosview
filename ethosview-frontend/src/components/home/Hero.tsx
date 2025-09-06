"use client";
import React from "react";
import type { DashboardResponse, AnalyticsSummaryResponse } from "../../types/api";
import { useCountUp } from "./useCountUp";

export function Hero({
  dashboard,
  analytics,
}: {
  dashboard: DashboardResponse;
  analytics: AnalyticsSummaryResponse;
}) {
  const totalCompanies = useCountUp(dashboard.summary.total_companies);
  const totalSectors = useCountUp(dashboard.summary.total_sectors);
  const avgESG = useCountUp(dashboard.summary.avg_esg_score, 800);
  const kpis = [
    { label: "Total Companies", value: Math.round(totalCompanies).toLocaleString() },
    { label: "Sectors", value: Math.round(totalSectors).toString() },
    { label: "Avg ESG", value: avgESG.toFixed(2) },
  ];

  return (
    <section className="py-10 md:py-14 relative overflow-hidden" style={{ background: "linear-gradient(180deg, #F5F7FA 0%, #E6F7F2 100%)" }}>
      <div className="blob blob-drift-a" style={{ top: -60, left: -80, width: 220, height: 220, background: "rgba(30,106,225,0.18)", borderRadius: 9999 }} />
      <div className="blob blob-drift-b" style={{ top: -40, right: -60, width: 200, height: 200, background: "rgba(42,179,166,0.18)", borderRadius: 9999 }} />
      <div className="max-w-6xl mx-auto px-4 relative">
        <h1 className="text-3xl md:text-5xl font-semibold tracking-tight text-gradient">EthosView</h1>
        <p className="mt-3 text-base md:text-lg" style={{ color: "#374151" }}>
          ESG and financial insights, unified. Live from our Go backend.
        </p>
        <div className="mt-3 flex items-center gap-3 text-xs" style={{ color: "#374151" }}>
          <span className="inline-flex items-center gap-2 glass-card px-2 py-1">
            <span style={{ width: 6, height: 6, borderRadius: 9999, background: "#1D9A6C", display: "inline-block" }} />
            Live
          </span>
          <span>Data updated just now</span>
        </div>
        <div className="mt-6 grid grid-cols-1 sm:grid-cols-3 gap-4 animate-fade-in-up">
          {kpis.map((kpi) => (
            <div key={kpi.label} className="glass-card p-4 hover-lift animate-fade-in-up" style={{ animationDuration: "360ms" }}>
              <div className="text-sm" style={{ color: "#6B7280" }}>{kpi.label}</div>
              <div className="mt-1 text-2xl font-medium" style={{ color: "#0B2545", fontVariantNumeric: "tabular-nums" }}>{kpi.value}</div>
            </div>
          ))}
        </div>
        <div className="mt-8 flex items-center gap-3">
          <a href="#market" className="btn-primary btn-sheen rounded px-4 py-2 hover-lift pulse-outline">View market</a>
          <a href="#featured" className="rounded px-4 py-2 tilt-hover" style={{ color: "#1E6AE1" }}>Featured →</a>
        </div>
        <div className="mt-10 flex justify-center float-slow">
          <a href="#market" className="text-gradient text-sm hover-lift">Scroll to explore ↓</a>
        </div>
      </div>
    </section>
  );
}


