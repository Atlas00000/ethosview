"use client";
import React, { useEffect, useMemo, useRef, useState } from "react";
import type { AnalyticsSummaryResponse, ESGTrendsResponse, LatestESGResponse, LatestPriceResponse, CompanyResponse } from "../../types/api";
import { api } from "../../services/api";
import { ResponsiveContainer, Tooltip, XAxis, YAxis, PieChart, Pie, Cell, BarChart, Bar, ReferenceLine } from "recharts";
import { useCountUp } from "./useCountUp";

type Props = { analytics: AnalyticsSummaryResponse };

export function ESGHighlightsPro({ analytics }: Props) {
  const list = analytics.top_esg_performers || [];
  const [selectedId, setSelectedId] = useState<number | null>(list.length ? list[0].company_id : null);
  const [trends, setTrends] = useState<ESGTrendsResponse | null>(null);
  const [latest, setLatest] = useState<LatestESGResponse | null>(null);
  const [price, setPrice] = useState<LatestPriceResponse | null>(null);
  const [company, setCompany] = useState<CompanyResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [days, setDays] = useState<number>(60);
  const [activeSlice, setActiveSlice] = useState<number | null>(null);
  const [paused, setPaused] = useState(false);
  const userInteractedRef = useRef(false);
  const formatNumber = (value: unknown, fractionDigits = 2): string => {
    return typeof value === 'number' && isFinite(value) ? value.toFixed(fractionDigits) : '—';
  };

  useEffect(() => {
    if (!selectedId) return;
    setLoading(true);
    Promise.all([
      api.esgTrends(selectedId, days),
      api.latestESG(selectedId),
      api.latestPrice(selectedId),
      api.companyById(selectedId),
    ])
      .then(([t, l, p, c]) => {
        setTrends(t);
        setLatest(l);
        setPrice(p);
        setCompany(c as CompanyResponse);
      })
      .finally(() => setLoading(false));
  }, [selectedId, days]);

  useEffect(() => {
    const idList = (list || []).slice(0, 5).map((x) => x.company_id);
    if (idList.length < 2) return;
    const timer = setInterval(() => {
      if (paused || userInteractedRef.current) return;
      setSelectedId((prev) => {
        if (!prev) return idList[0];
        const idx = idList.indexOf(prev);
        const next = idx >= 0 && idx < idList.length - 1 ? idList[idx + 1] : idList[0];
        return next;
      });
    }, 8000);
    return () => clearInterval(timer);
  }, [list, paused]);

  const chartData = useMemo(() => {
    const rows = Array.isArray(trends?.trends) ? trends!.trends.slice().reverse() : [];
    return rows.map((r) => ({ x: r.date, y: r.esg_score, e: r.e_score, s: r.s_score, g: r.g_score }));
  }, [trends]);

  const latestPoint = chartData.length ? chartData[chartData.length - 1] : null;
  const prevPoint = chartData.length > 1 ? chartData[chartData.length - 2] : null;
  const pieData = latestPoint
    ? [
        { name: "E", value: latestPoint.e },
        { name: "S", value: latestPoint.s },
        { name: "G", value: latestPoint.g },
      ]
    : [];
  const overallValue = typeof latest?.overall_score === 'number' ? latest.overall_score : (latestPoint ? latestPoint.y : undefined);
  const overallDelta = latestPoint && prevPoint ? latestPoint.y - prevPoint.y : 0;
  const overallDeltaPct = latestPoint && prevPoint && prevPoint.y ? ((latestPoint.y - prevPoint.y) / prevPoint.y) * 100 : 0;
  const barData = latestPoint && typeof overallValue === 'number'
    ? [
        { name: "E", score: latestPoint.e },
        { name: "S", score: latestPoint.s },
        { name: "G", score: latestPoint.g },
        { name: "ESG", score: overallValue },
      ]
    : [];

  const countE = useCountUp(typeof latestPoint?.e === 'number' ? latestPoint!.e : 0);
  const countS = useCountUp(typeof latestPoint?.s === 'number' ? latestPoint!.s : 0);
  const countG = useCountUp(typeof latestPoint?.g === 'number' ? latestPoint!.g : 0);
  const countESG = useCountUp(typeof overallValue === 'number' ? overallValue : 0);

  // Sector average comparison for overall ESG (when company sector available)
  const sectorAvg = useMemo(() => {
    if (!company?.sector) return null as number | null;
    const row = analytics.sector_comparisons.find(s => s.sector === company.sector);
    return row ? row.avg_esg_score : null;
  }, [company?.sector, analytics.sector_comparisons]);

  return (
    <section className="max-w-6xl mx-auto px-4 py-12 relative overflow-hidden"
      onMouseEnter={() => setPaused(true)}
      onMouseLeave={() => setPaused(false)}
    >
      <h2 className="text-xl font-semibold mb-3 text-gradient">ESG highlights</h2>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="md:col-span-1 glass-card divide-y animate-fade-in-up">
          {(list || []).slice(0, 15).map((it) => (
            <button
              key={it.company_id}
              onClick={() => { setSelectedId(it.company_id); userInteractedRef.current = true; }}
              className={`w-full text-left p-3 hover-lift ${selectedId === it.company_id ? 'bg-white/40' : ''}`}
            >
              <div className="text-sm font-medium" style={{ color: "#0B2545" }}>{it.company_name}</div>
              <div className="text-xs" style={{ color: "#374151" }}>ESG: {formatNumber(it.value, 2)} / 100 • Rank #{it.rank ?? '—'}</div>
            </button>
          ))}
        </div>
        <div className="md:col-span-2 glass-card p-4 animate-fade-in-up">
          <div className="flex flex-wrap items-end justify-between gap-3">
            <div>
              <div className="text-xs" style={{ color: "#374151" }}>Company</div>
              <div className="text-lg font-semibold" style={{ color: "#0B2545" }}>
                {list.find((x) => x.company_id === selectedId)?.company_name || '—'}
              </div>
            </div>
            <div className="text-right">
              <div className="text-xs" style={{ color: "#374151" }}>Latest ESG</div>
              <div className="text-lg font-semibold" style={{ color: "#1D9A6C" }}>{formatNumber(latest?.overall_score ?? latestPoint?.y, 2)} / 100</div>
              {latestPoint && (
                <div className="text-xs" style={{ color: overallDelta >= 0 ? '#1D9A6C' : '#E25555' }}>
                  {(overallDelta >= 0 ? '▲ +' : '▼ ') + formatNumber(overallDelta, 2)} ({formatNumber(overallDeltaPct, 2)}%)
                </div>
              )}
            </div>
            <div className="text-right">
              <div className="text-xs" style={{ color: "#374151" }}>Price</div>
              <div className="text-lg font-semibold" style={{ color: "#0B2545" }}>{price?.price?.close_price?.toFixed ? price.price.close_price.toFixed(2) : '—'}</div>
            </div>
          </div>

          <div className="mt-3 grid grid-cols-4 gap-2">
            <div className="glass-card p-2 text-center">
              <div className="text-[10px]" style={{ color: "#374151" }}>E</div>
              <div className="text-base font-semibold" style={{ color: "#2AB3A6" }}>{formatNumber(countE, 1)}</div>
            </div>
            <div className="glass-card p-2 text-center">
              <div className="text-[10px]" style={{ color: "#374151" }}>S</div>
              <div className="text-base font-semibold" style={{ color: "#1E6AE1" }}>{formatNumber(countS, 1)}</div>
            </div>
            <div className="glass-card p-2 text-center">
              <div className="text-[10px]" style={{ color: "#374151" }}>G</div>
              <div className="text-base font-semibold" style={{ color: "#9AD36A" }}>{formatNumber(countG, 1)}</div>
            </div>
            <div className="glass-card p-2 text-center">
              <div className="text-[10px]" style={{ color: "#374151" }}>ESG</div>
              <div className="text-base font-semibold" style={{ color: "#1D9A6C" }}>{formatNumber(countESG, 1)}</div>
            </div>
          </div>

          <div className="mt-3 flex items-center gap-2">
            {[30, 60, 180].map((d) => (
              <button key={d} onClick={() => { setDays(d); userInteractedRef.current = true; }}
                className={`px-3 py-1 rounded-full text-xs transition-all ${days === d ? 'bg-white/60 shadow-elevated' : 'bg-white/30 hover:bg-white/50'} tilt-hover`}
              >
                {d}d
              </button>
            ))}
          </div>

          <div className="mt-3 grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div className="glass-card p-3 tilt-hover">
              <div className="text-xs mb-1" style={{ color: "#374151" }}>ESG component breakdown</div>
              <div style={{ width: '100%', height: 180 }}>
                {pieData.length ? (
                  <ResponsiveContainer>
                    <PieChart>
                      <Pie data={pieData} dataKey="value" nameKey="name" innerRadius={36} outerRadius={60} paddingAngle={2}
                        onMouseEnter={(_, idx) => setActiveSlice(idx)} onMouseLeave={() => setActiveSlice(null)}
                      >
                        <Cell fill="#2AB3A6" {...(activeSlice === 0 ? { outerRadius: 66 } : {})} />
                        <Cell fill="#1E6AE1" {...(activeSlice === 1 ? { outerRadius: 66 } : {})} />
                        <Cell fill="#9AD36A" {...(activeSlice === 2 ? { outerRadius: 66 } : {})} />
                      </Pie>
                      <Tooltip formatter={(v: any, n: any) => [formatNumber(v, 2), n]} />
                    </PieChart>
                  </ResponsiveContainer>
                ) : (
                  <div className="skeleton w-full h-full rounded" />
                )}
              </div>
              <div className="mt-2 flex items-center gap-2 text-xs" style={{ color: '#374151' }}>
                <span className="px-2 py-0.5 rounded-full bg-white/50">E = Environmental</span>
                <span className="px-2 py-0.5 rounded-full bg-white/50">S = Social</span>
                <span className="px-2 py-0.5 rounded-full bg-white/50">G = Governance</span>
              </div>
            </div>
            <div className="glass-card p-3 tilt-hover">
              <div className="text-xs mb-1" style={{ color: "#374151" }}>Components vs overall (score)</div>
              <div style={{ width: '100%', height: 180 }}>
                {barData.length ? (
                  <ResponsiveContainer>
                    <BarChart data={barData} margin={{ top: 6, right: 8, left: 0, bottom: 0 }}>
                      <XAxis dataKey="name" tick={{ fontSize: 10 }} />
                      <YAxis domain={[0, 100]} tick={{ fontSize: 10 }} />
                      <Tooltip formatter={(v: any, n: any) => [formatNumber(v, 2), n]} />
                      {sectorAvg != null && <ReferenceLine y={sectorAvg} stroke="#1D9A6C" strokeDasharray="4 4" label={{ value: 'Sector Avg', position: 'insideTopRight', fill: '#1D9A6C', fontSize: 10 }} />}
                      <Bar dataKey="score" radius={[6, 6, 0, 0]} fill="#0B2545" className="hover-lift" isAnimationActive />
                    </BarChart>
                  </ResponsiveContainer>
                ) : (
                  <div className="skeleton w-full h-full rounded" />
                )}
              </div>
            </div>
          </div>

          {loading && <div className="mt-2 hero-progress"><span /></div>}
        </div>
      </div>
    </section>
  );
}


