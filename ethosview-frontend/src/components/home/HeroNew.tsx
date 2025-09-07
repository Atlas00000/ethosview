"use client";
import React, { useEffect, useMemo, useRef, useState } from "react";
import type { DashboardResponse, AnalyticsSummaryResponse, MarketLatestResponse, MarketHistoryResponse, ESGTrendsResponse } from "../../types/api";
import { api } from "../../services/api";
import { useCountUp } from "./useCountUp";
import { ResponsiveContainer, AreaChart, Area, Tooltip, XAxis, YAxis, PieChart, Pie, Cell, BarChart, Bar, RadialBarChart, RadialBar } from "recharts";

type Props = {
  dashboard: DashboardResponse;
  analytics: AnalyticsSummaryResponse;
  market: MarketLatestResponse;
  history: MarketHistoryResponse | null;
};

export function HeroNew({ dashboard, analytics, market, history }: Props) {
  const [index, setIndex] = useState(0);
  const [paused, setPaused] = useState(false);
  const containerRef = useRef<HTMLDivElement | null>(null);

  const topESG = (analytics?.top_esg_performers || []).slice(0, 1);
  const initialCompanyId = topESG.length ? topESG[0].company_id : null;
  const [companyId, setCompanyId] = useState<number | null>(initialCompanyId);
  const [trends, setTrends] = useState<ESGTrendsResponse | null>(null);

  useEffect(() => {
    if (!companyId) return;
    api.esgTrends(companyId, 60).then(setTrends).catch(() => setTrends(null));
  }, [companyId]);

  const slideCount = 4;
  useEffect(() => {
    const id = setInterval(() => {
      if (!paused) setIndex((i) => (i + 1) % slideCount);
    }, 7000);
    return () => clearInterval(id);
  }, [paused]);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;
    function onKey(e: KeyboardEvent) {
      if (e.key === "ArrowRight") setIndex((i) => (i + 1) % slideCount);
      if (e.key === "ArrowLeft") setIndex((i) => (i + slideCount - 1) % slideCount);
    }
    el.addEventListener("keydown", onKey as any);
    return () => el.removeEventListener("keydown", onKey as any);
  }, []);

  const totalCompanies = useCountUp(dashboard.summary.total_companies);
  const totalSectors = useCountUp(dashboard.summary.total_sectors);
  const avgESG = useCountUp(dashboard.summary.avg_esg_score, 800);

  const sectors = useMemo(() => {
    const rows = (analytics?.sector_comparisons || []).slice();
    rows.sort((a, b) => b.avg_esg_score - a.avg_esg_score);
    return rows.slice(0, 6);
  }, [analytics?.sector_comparisons]);

  const marketData = market.market_data;
  const ticker = [
    { label: "S&P 500", value: marketData.sp500_close },
    { label: "NASDAQ", value: marketData.nasdaq_close },
    { label: "DOW", value: marketData.dow_close },
    { label: "VIX", value: marketData.vix_close },
  ].filter((x) => typeof x.value === "number");

  const spSeries = useMemo(() => {
    const rows = Array.isArray(history?.data) ? history!.data.slice().reverse() : [];
    return rows.map((d) => ({ x: d.date, y: d.sp500_close ?? d.nasdaq_close ?? d.dow_close ?? 0 }));
  }, [history]);

  const esgSeries = useMemo(() => {
    const rows = Array.isArray(trends?.trends) ? trends!.trends.slice().reverse() : [];
    return rows.map((t) => ({ x: t.date, y: t.esg_score }));
  }, [trends]);

  const topCap = (analytics?.top_market_cap || []).slice(0, 8);

  // Additional visuals
  const sectorPie = useMemo(() => {
    const rows = (analytics?.sector_comparisons || []).slice();
    rows.sort((a, b) => b.company_count - a.company_count);
    const top = rows.slice(0, 6);
    const otherCount = rows.slice(6).reduce((acc, r) => acc + (r.company_count || 0), 0);
    const data = top.map((r) => ({ name: r.sector, value: r.company_count }));
    if (otherCount > 0) data.push({ name: "Other", value: otherCount });
    return data;
  }, [analytics?.sector_comparisons]);

  const sectorESGBars = useMemo(() => {
    const rows = (analytics?.sector_comparisons || []).slice();
    rows.sort((a, b) => b.avg_esg_score - a.avg_esg_score);
    return rows.slice(0, 6).map((r) => ({ name: r.sector, esg: r.avg_esg_score }));
  }, [analytics?.sector_comparisons]);

  const avgEsgValue = Math.max(0, Math.min(100, dashboard.summary.avg_esg_score || 0));

  // Compute padded Y domains to avoid visually flat charts
  function paddedDomain(series: Array<{ y: number }>, padPercent = 0.06, minClamp?: number, maxClamp?: number): [number, number] {
    if (!series.length) return [0, 1];
    let min = series.reduce((m, p) => Math.min(m, Number(p.y) || 0), Number.POSITIVE_INFINITY);
    let max = series.reduce((m, p) => Math.max(m, Number(p.y) || 0), Number.NEGATIVE_INFINITY);
    if (!isFinite(min) || !isFinite(max)) return [0, 1];
    if (min === max) {
      const delta = Math.max(1, Math.abs(min) * 0.05);
      min -= delta;
      max += delta;
    }
    const pad = (max - min) * padPercent;
    min -= pad;
    max += pad;
    if (typeof minClamp === 'number') min = Math.max(min, minClamp);
    if (typeof maxClamp === 'number') max = Math.min(max, maxClamp);
    return [min, max];
  }
  const spDomain = paddedDomain(spSeries);
  const esgDomain = paddedDomain(esgSeries, 0.08, 0, 100);

  return (
    <section className="py-10 md:py-14 relative overflow-hidden" style={{ background: "linear-gradient(180deg, #F5F7FA 0%, #E6F7F2 100%)" }}>
      <div className="ribbon ribbon-blue" style={{ top: 0, left: "-10%" }} />
      <div className="ribbon ribbon-green" style={{ bottom: -20, left: "-5%", ['--rot' as any]: '5deg' }} />
      <div className="blob blob-drift-a" style={{ top: -60, left: -80, width: 220, height: 220, background: "rgba(30,106,225,0.18)", borderRadius: 9999 }} />
      <div className="blob blob-drift-b" style={{ top: -40, right: -60, width: 200, height: 200, background: "rgba(42,179,166,0.18)", borderRadius: 9999 }} />
      <div className="max-w-6xl mx-auto px-4 relative" ref={containerRef} tabIndex={0} onMouseEnter={() => setPaused(true)} onMouseLeave={() => setPaused(false)}>
        <h1 className="text-3xl md:text-5xl font-semibold tracking-tight text-gradient">EthosView</h1>
        <p className="mt-3 text-base md:text-lg" style={{ color: "#374151" }}>Welcome to EthosView, your real time ESG and financial intelligence hub.</p>

        {ticker.length > 0 && (
          <div className="mt-4 flex gap-2 overflow-auto scrollbar-hide">
            {ticker.map((t) => (
              <span key={t.label} className="glass-card px-3 py-1 text-xs hover-lift btn-sheen" style={{ color: "#0B2545" }}>
                {t.label}: {Number(t.value).toLocaleString(undefined, { maximumFractionDigits: 2 })}
              </span>
            ))}
          </div>
        )}

        <div className="relative mt-6 min-h-[240px]">
          <div className={`slide-base ${index === 0 ? "slide-active" : ""}`}>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 animate-fade-in-up">
              <div className="glass-card p-4 hover-lift spotlight-hover card-glow">
                <div className="text-sm" style={{ color: "#6B7280" }}>Total Companies</div>
                <div className="mt-1 text-2xl font-medium" style={{ color: "#0B2545", fontVariantNumeric: "tabular-nums" }}>{Math.round(totalCompanies).toLocaleString()}</div>
              </div>
              <div className="glass-card p-4 hover-lift spotlight-hover card-glow">
                <div className="text-sm" style={{ color: "#6B7280" }}>Sectors</div>
                <div className="mt-1 text-2xl font-medium" style={{ color: "#0B2545" }}>{Math.round(totalSectors)}</div>
              </div>
              <div className="glass-card p-4 hover-lift spotlight-hover card-glow">
                <div className="text-sm" style={{ color: "#6B7280" }}>Avg ESG</div>
                <div className="mt-1 text-2xl font-medium" style={{ color: "#1D9A6C" }}>{avgESG.toFixed(2)} / 100</div>
              </div>
            </div>
            {sectors.length > 0 && (
              <div className="mt-4 grid grid-cols-2 sm:grid-cols-3 gap-3">
                {sectors.map((s) => (
                  <div key={s.sector} className="glass-card p-3 hover-lift tilt-hover">
                    <div className="text-xs" style={{ color: "#374151" }}>{s.sector}</div>
                    <div className="mt-1 text-sm font-medium" style={{ color: "#0B2545" }}>Avg ESG {s.avg_esg_score.toFixed(2)}</div>
                    <div className="mt-2 h-1.5 bg-white/50 rounded">
                      <div className="h-1.5 rounded" style={{ width: `${Math.min(100, s.avg_esg_score)}%`, background: "linear-gradient(90deg,#1E6AE1,#2AB3A6)" }} />
                    </div>
                    <div className="mt-1 text-[10px]" style={{ color: "#6B7280" }}>{s.company_count} companies</div>
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className={`slide-base ${index === 1 ? "slide-active" : ""}`}>
            <div className="animate-fade-in-up">
              <div className="text-sm mb-2" style={{ color: "#374151" }}>Top by market cap</div>
              <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
                {topCap.map((c) => (
                  <div key={c.company_id} className="glass-card p-3 hover-lift tilt-hover">
                    <div className="text-xs" style={{ color: "#6B7280" }}>{c.company_name}</div>
                    <div className="mt-1 text-base font-semibold" style={{ color: "#0B2545" }}>
                      {Number(c.value).toLocaleString(undefined, { notation: "compact", maximumFractionDigits: 2 })}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className={`slide-base ${index === 2 ? "slide-active" : ""}`}>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 animate-fade-in-up">
              <div className="glass-card p-3">
                <div className="text-xs" style={{ color: "#374151" }}>S&P 500 (pts)</div>
                <div style={{ width: "100%", height: 120 }}>
                  <ResponsiveContainer>
                    <BarChart data={spSeries} margin={{ top: 4, right: 8, left: 0, bottom: 0 }}>
                      <defs>
                        <linearGradient id="heroSpBar" x1="0" y1="0" x2="0" y2="1">
                          <stop offset="10%" stopColor="#1E6AE1" stopOpacity={0.9} />
                          <stop offset="90%" stopColor="#1E6AE1" stopOpacity={0.2} />
                        </linearGradient>
                      </defs>
                      <XAxis dataKey="x" hide />
                      <YAxis hide domain={spDomain as any} />
                      <Tooltip formatter={(v: any) => [Number(v).toLocaleString(undefined, { maximumFractionDigits: 2 }), "Index"]} />
                      <Bar dataKey="y" fill="url(#heroSpBar)" radius={[4,4,0,0]} barSize={6} isAnimationActive />
                    </BarChart>
                  </ResponsiveContainer>
                </div>
              </div>
              <div className="glass-card p-3">
                <div className="text-xs" style={{ color: "#374151" }}>ESG trend (top performer)</div>
                <div style={{ width: "100%", height: 120 }}>
                  <ResponsiveContainer>
                    <BarChart data={esgSeries} margin={{ top: 4, right: 8, left: 0, bottom: 0 }}>
                      <defs>
                        <linearGradient id="heroEsgBar" x1="0" y1="0" x2="0" y2="1">
                          <stop offset="10%" stopColor="#1D9A6C" stopOpacity={0.9} />
                          <stop offset="90%" stopColor="#1D9A6C" stopOpacity={0.2} />
                        </linearGradient>
                      </defs>
                      <XAxis dataKey="x" hide />
                      <YAxis hide domain={esgDomain as any} />
                      <Tooltip formatter={(v: any) => [Number(v).toFixed(2), "Score"]} />
                      <Bar dataKey="y" fill="url(#heroEsgBar)" radius={[4,4,0,0]} barSize={6} isAnimationActive />
                    </BarChart>
                  </ResponsiveContainer>
                </div>
              </div>
            </div>
          </div>

          <div className={`slide-base ${index === 3 ? "slide-active" : ""}`}>
            <div className="mt-0 grid grid-cols-1 md:grid-cols-3 gap-4 animate-fade-in-up">
              <div className="glass-card p-3 tilt-hover">
                <div className="text-xs mb-1" style={{ color: "#374151" }}>Sector distribution (companies)</div>
                <div style={{ width: "100%", height: 160 }}>
                  {sectorPie.length ? (
                    <ResponsiveContainer>
                      <PieChart>
                        <Pie data={sectorPie} dataKey="value" nameKey="name" innerRadius={40} outerRadius={68} paddingAngle={2}>
                          {sectorPie.map((_, idx) => (
                            <Cell key={idx} fill={["#1E6AE1","#2AB3A6","#9AD36A","#0B2545","#7C3AED","#EC4899","#F59E0B"][idx % 7]} />
                          ))}
                        </Pie>
                        <Tooltip formatter={(v: any, n: any) => [Number(v).toLocaleString(), n]} />
                      </PieChart>
                    </ResponsiveContainer>
                  ) : (
                    <div className="skeleton w-full h-full rounded" />
                  )}
                </div>
              </div>

              <div className="glass-card p-3 tilt-hover">
                <div className="text-xs mb-1" style={{ color: "#374151" }}>Top sectors by avg ESG</div>
                <div style={{ width: "100%", height: 160 }}>
                  {sectorESGBars.length ? (
                    <ResponsiveContainer>
                      <BarChart data={sectorESGBars} margin={{ top: 8, right: 8, left: 0, bottom: 0 }}>
                        <XAxis dataKey="name" tick={{ fontSize: 10 }} interval={0} angle={-15} height={30} />
                        <YAxis domain={[0, 100]} tick={{ fontSize: 10 }} />
                        <Tooltip formatter={(v: any) => [Number(v).toFixed(2), "Avg ESG"]} />
                        <Bar dataKey="esg" radius={[6,6,0,0]} fill="#1D9A6C" />
                      </BarChart>
                    </ResponsiveContainer>
                  ) : (
                    <div className="skeleton w-full h-full rounded" />
                  )}
                </div>
              </div>

              <div className="glass-card p-3 tilt-hover flex items-center justify-center">
                <div style={{ width: "100%,", height: 160 }}>
                  <ResponsiveContainer>
                    <RadialBarChart innerRadius="60%" outerRadius="90%" data={[{ name: "Avg ESG", value: avgEsgValue }] } startAngle={90} endAngle={-270}>
                      <RadialBar minAngle={15} background clockWise dataKey="value" cornerRadius={8} fill="#1D9A6C" />
                      <Tooltip formatter={(v: any) => [Number(v).toFixed(2), "Avg ESG"]} />
                    </RadialBarChart>
                  </ResponsiveContainer>
                  <div className="-mt-24 text-center">
                    <div className="text-xs" style={{ color: "#374151" }}>Global Avg ESG</div>
                    <div className="text-xl font-semibold" style={{ color: "#0B2545" }}>{avgEsgValue.toFixed(2)} / 100</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div className="hero-progress mt-4"><span /></div>
        </div>

        <div className="mt-6 flex items-center gap-2">
          {Array.from({ length: slideCount }, (_, i) => i).map((i) => (
            <button key={i} aria-label={`Slide ${i + 1}`} onClick={() => setIndex(i)} className="rounded-full" style={{ width: 8, height: 8, background: i === index ? "#1E6AE1" : "#D1D5DB" }} />
          ))}
        </div>

        
      </div>
    </section>
  );
}


