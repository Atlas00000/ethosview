"use client";
import React, { useEffect, useMemo, useState } from "react";
import type { DashboardResponse, AnalyticsSummaryResponse, MarketLatestResponse, MarketHistoryResponse, ESGTrendsResponse } from "../../types/api";
import { useCountUp } from "./useCountUp";
import { ResponsiveContainer, LineChart, Line } from "recharts";

type Props = {
  dashboard: DashboardResponse;
  analytics: AnalyticsSummaryResponse;
  market: MarketLatestResponse;
  history?: MarketHistoryResponse | null;
  esgTrends?: ESGTrendsResponse | null;
};

export function HeroSlider({ dashboard, analytics, market, history, esgTrends }: Props) {
  const [index, setIndex] = useState(0);
  const slides = useMemo(() => [
    <KPISlide key="kpi" dashboard={dashboard} />,
    <MarketSlide key="market" market={market} />,
    <ESGSlide key="esg" analytics={analytics} />,
    <ChartsSlide key="charts" history={history || undefined} esgTrends={esgTrends || undefined} />,
  ], [dashboard, analytics, market, history, esgTrends]);

  useEffect(() => {
    const id = setInterval(() => setIndex((i) => (i + 1) % slides.length), 6000);
    return () => clearInterval(id);
  }, [slides.length]);

  return (
    <section className="py-10 md:py-14 relative overflow-hidden" style={{ background: "linear-gradient(180deg, #F5F7FA 0%, #E6F7F2 100%)" }}>
      <div className="blob blob-drift-a" style={{ top: -60, left: -80, width: 220, height: 220, background: "rgba(30,106,225,0.18)", borderRadius: 9999 }} />
      <div className="blob blob-drift-b" style={{ top: -40, right: -60, width: 200, height: 200, background: "rgba(42,179,166,0.18)", borderRadius: 9999 }} />
      <div className="max-w-6xl mx-auto px-4 relative">
        <h1 className="text-3xl md:text-5xl font-semibold tracking-tight text-gradient">EthosView</h1>
        <p className="mt-3 text-base md:text-lg" style={{ color: "#374151" }}>
          ESG and financial insights, unified. Live from our Go backend.
        </p>

        <div className="relative mt-6 min-h-[200px]" style={{ position: "relative", ['--slide-dur' as any]: '6000ms' }}>
          {slides.map((node, i) => (
            <div key={i} className={`slide-base ${i === index ? 'slide-active' : ''}`}>{node}</div>
          ))}
          <div className="hero-progress mt-4"><span /></div>
        </div>

        <div className="mt-6 flex items-center gap-2">
          {slides.map((_, i) => (
            <button
              key={i}
              aria-label={`Slide ${i + 1}`}
              onClick={() => setIndex(i)}
              className="rounded-full"
              style={{ width: 8, height: 8, background: i === index ? "#1E6AE1" : "#D1D5DB" }}
            />
          ))}
        </div>

        <div className="mt-8 flex items-center gap-3">
          <a href="#market" className="btn-primary btn-sheen rounded px-4 py-2 hover-lift pulse-outline">View market</a>
          <a href="#featured" className="rounded px-4 py-2 tilt-hover" style={{ color: "#1E6AE1" }}>Featured →</a>
        </div>
      </div>
    </section>
  );
}

function KPISlide({ dashboard }: { dashboard: DashboardResponse }) {
  const totalCompanies = useCountUp(dashboard.summary.total_companies);
  const totalSectors = useCountUp(dashboard.summary.total_sectors);
  const avgESG = useCountUp(dashboard.summary.avg_esg_score, 800);
  const kpis = [
    { label: "Total Companies", value: `${Math.round(totalCompanies).toLocaleString()} companies` },
    { label: "Sectors", value: `${Math.round(totalSectors).toString()} sectors` },
    { label: "Avg ESG", value: `${avgESG.toFixed(2)} / 100` },
  ];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 animate-fade-in-up">
      {kpis.map((kpi) => (
        <div key={kpi.label} className="glass-card p-4 hover-lift">
          <div className="text-sm" style={{ color: "#6B7280" }}>{kpi.label}</div>
          <div className="mt-1 text-2xl font-medium" style={{ color: "#0B2545", fontVariantNumeric: "tabular-nums" }}>{kpi.value}</div>
        </div>
      ))}
    </div>
  );
}

function MarketSlide({ market }: { market: MarketLatestResponse }) {
  const m = market.market_data;
  const items = [
    { label: "S&P 500", value: m.sp500_close },
    { label: "NASDAQ", value: m.nasdaq_close },
    { label: "DOW", value: m.dow_close },
    { label: "VIX", value: m.vix_close },
  ];
  return (
    <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 animate-fade-in-up">
      {items.map((it) => (
        <div key={it.label} className="glass-card p-4 hover-lift">
          <div className="text-sm" style={{ color: "#6B7280" }}>{it.label}</div>
          <div className="mt-1 text-2xl font-medium" style={{ color: "#0B2545" }}>
            {typeof it.value === "number" ? `${it.value.toLocaleString(undefined, { maximumFractionDigits: 2 })} pts` : "—"}
          </div>
        </div>
      ))}
    </div>
  );
}

function ESGSlide({ analytics }: { analytics: AnalyticsSummaryResponse }) {
  const top = (analytics?.top_esg_performers || []).slice(0, 3);
  return (
    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 animate-fade-in-up">
      {top.map((t) => (
        <div key={t.company_id} className="glass-card p-4 hover-lift">
          <div className="text-sm" style={{ color: "#6B7280" }}>{t.company_name}</div>
          <div className="mt-1 text-2xl font-medium" style={{ color: "#1D9A6C" }}>{(Number.isFinite(t.value) ? t.value : 0).toFixed(2)} / 100</div>
        </div>
      ))}
    </div>
  );
}

function ChartsSlide({ history, esgTrends }: { history?: MarketHistoryResponse; esgTrends?: ESGTrendsResponse }) {
  const spData = useMemo(() => {
    const rows = Array.isArray(history?.data) ? history!.data.slice().reverse() : [];
    return rows.map((d) => ({ x: d.date, y: d.sp500_close ?? d.nasdaq_close ?? d.dow_close ?? 0 }));
  }, [history]);
  const esgData = useMemo(() => {
    const rows = Array.isArray(esgTrends?.trends) ? esgTrends!.trends.slice().reverse() : [];
    return rows.map((t) => ({ x: t.date, y: t.esg_score }));
  }, [esgTrends]);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div className="glass-card p-3">
        <div className="text-xs" style={{ color: "#374151" }}>S&amp;P 500 (pts)</div>
        <div style={{ width: "100%", height: 90 }}>
          <ResponsiveContainer>
            <LineChart data={spData} margin={{ top: 2, bottom: 0, left: 0, right: 0 }}>
              <Line type="monotone" dataKey="y" stroke="#1E6AE1" strokeWidth={1.5} dot={false} />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>
      <div className="glass-card p-3">
        <div className="text-xs" style={{ color: "#374151" }}>ESG trend (score)</div>
        <div style={{ width: "100%", height: 90 }}>
          <ResponsiveContainer>
            <LineChart data={esgData} margin={{ top: 2, bottom: 0, left: 0, right: 0 }}>
              <Line type="monotone" dataKey="y" stroke="#1D9A6C" strokeWidth={1.5} dot={false} />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
}


