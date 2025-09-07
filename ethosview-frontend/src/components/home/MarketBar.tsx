"use client";
import React, { useEffect, useMemo, useState } from "react";
import type { MarketLatestResponse, MarketHistoryResponse } from "../../types/api";
import { api } from "../../services/api";
import { ResponsiveContainer, AreaChart, Area, Tooltip, XAxis, YAxis, BarChart, Bar } from "recharts";
import { useCountUp } from "./useCountUp";

type RangeKey = "1d" | "1w" | "1m";

export function MarketBar({ market }: { market: MarketLatestResponse }) {
  const m = market.market_data;
  const [range, setRange] = useState<RangeKey>("1w");
  const [history, setHistory] = useState<MarketHistoryResponse | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const today = new Date();
    const days = range === "1d" ? 1 : range === "1w" ? 7 : 30;
    const start = new Date(today.getTime() - days * 24 * 3600 * 1000);
    const fmt = (d: Date) => d.toISOString().slice(0, 10);
    setLoading(true);
    api
      .marketHistory(fmt(start), fmt(today), 100)
      .then((res) => setHistory(res))
      .finally(() => setLoading(false));
  }, [range]);

  const seriesSP = useMemo(() => {
    const rows = Array.isArray(history?.data) ? history!.data : [];
    return rows
      .slice()
      .reverse()
      .map((d) => ({ x: d.date, y: d.sp500_close ?? d.nasdaq_close ?? d.dow_close ?? 0 }));
  }, [history]);

  const seriesNDQ = useMemo(() => {
    const rows = Array.isArray(history?.data) ? history!.data : [];
    return rows
      .slice()
      .reverse()
      .map((d) => ({ x: d.date, y: d.nasdaq_close ?? d.sp500_close ?? d.dow_close ?? 0 }));
  }, [history]);

  const seriesDOW = useMemo(() => {
    const rows = Array.isArray(history?.data) ? history!.data : [];
    return rows
      .slice()
      .reverse()
      .map((d) => ({ x: d.date, y: d.dow_close ?? d.sp500_close ?? d.nasdaq_close ?? 0 }));
  }, [history]);

  const seriesVIX = useMemo(() => {
    const rows = Array.isArray(history?.data) ? history!.data : [];
    return rows
      .slice()
      .reverse()
      .map((d) => ({ x: d.date, y: (d.vix_close ?? m.vix_close ?? 0) as number }));
  }, [history]);

  const seriesY10 = useMemo(() => {
    const rows = Array.isArray(history?.data) ? history!.data : [];
    return rows
      .slice()
      .reverse()
      .map((d) => ({ x: d.date, y: (d.treasury_10y ?? m.treasury_10y ?? 0) as number }));
  }, [history]);

  function computeDelta(series: Array<{ x: string; y: number }>, fallback: number) {
    const last = series.length ? series[series.length - 1].y : fallback;
    const prev = series.length > 1 ? series[series.length - 2].y : last;
    const delta = last - prev;
    const pct = prev ? (delta / prev) * 100 : 0;
    return { last, prev, delta, pct };
  }

  const spDelta = computeDelta(seriesSP, m.sp500_close ?? 0);
  const ndqDelta = computeDelta(seriesNDQ, m.nasdaq_close ?? 0);
  const dowDelta = computeDelta(seriesDOW, m.dow_close ?? 0);

  // VIX and 10Y deltas from last two rows if present
  const rows = Array.isArray(history?.data) ? history!.data.slice().reverse() : [];
  const vixLast = rows.length ? (rows[rows.length - 1].vix_close ?? m.vix_close ?? 0) : (m.vix_close ?? 0);
  const vixPrev = rows.length > 1 ? (rows[rows.length - 2].vix_close ?? vixLast) : vixLast;
  const vixDelta = vixLast - vixPrev;
  const vixPct = vixPrev ? (vixDelta / vixPrev) * 100 : 0;

  const y10Last = rows.length ? (rows[rows.length - 1].treasury_10y ?? m.treasury_10y ?? 0) : (m.treasury_10y ?? 0);
  const y10Prev = rows.length > 1 ? (rows[rows.length - 2].treasury_10y ?? y10Last) : y10Last;
  const y10Delta = y10Last - y10Prev;
  const y10Pct = y10Prev ? (y10Delta / y10Prev) * 100 : 0;
  const trendColor = spDelta.delta > 0 ? "#1D9A6C" : spDelta.delta < 0 ? "#E25555" : "#374151";
  const trendArrow = spDelta.delta > 0 ? "▲" : spDelta.delta < 0 ? "▼" : "";

  // Avoid flat charts by padding domains
  function paddedDomain(series: Array<{ y: number }>, padPct = 0.06): [number, number] {
    if (!series.length) return [0, 1];
    let min = Number.POSITIVE_INFINITY;
    let max = Number.NEGATIVE_INFINITY;
    for (const p of series) {
      const y = Number(p.y) || 0;
      if (y < min) min = y;
      if (y > max) max = y;
    }
    if (!isFinite(min) || !isFinite(max)) return [0, 1];
    if (min === max) {
      const d = Math.max(1, Math.abs(min) * 0.05);
      min -= d;
      max += d;
    }
    const pad = (max - min) * padPct;
    return [min - pad, max + pad];
  }
  const trim = (arr: Array<{ x: string; y: number }>, n = 30) => arr.slice(Math.max(0, arr.length - n));
  const spBars = trim(seriesSP, 30);
  const ndqBars = trim(seriesNDQ, 30);
  const dowBars = trim(seriesDOW, 30);
  const vixBars = trim(seriesVIX, 30);
  const y10Bars = trim(seriesY10, 30);
  const spDomain = paddedDomain(spBars);
  const ndqDomain = paddedDomain(ndqBars);
  const dowDomain = paddedDomain(dowBars);

  // Range position for S&P
  const spMin = seriesSP.length ? Math.min(...seriesSP.map(s => s.y)) : 0;
  const spMax = seriesSP.length ? Math.max(...seriesSP.map(s => s.y)) : 1;
  const spPosPct = spMax !== spMin ? ((spDelta.last - spMin) / (spMax - spMin)) * 100 : 50;

  const sp = useCountUp(m.sp500_close ?? 0);
  const ndq = useCountUp(m.nasdaq_close ?? 0);
  const dow = useCountUp(m.dow_close ?? 0);
  const vix = useCountUp(m.vix_close ?? 0);
  const y10 = useCountUp(m.treasury_10y ?? 0);

  return (
    <section className="py-3">
      <div className="max-w-6xl mx-auto px-4">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div className="flex flex-wrap items-center gap-4 text-sm" style={{ color: "#0B2545" }}>
            <div className="font-medium">Market Snapshot</div>
            <div className="opacity-80">S&P 500: {sp.toLocaleString(undefined, { maximumFractionDigits: 2 })} pts</div>
            <div className="opacity-80">NASDAQ: {ndq.toLocaleString(undefined, { maximumFractionDigits: 2 })} pts</div>
            <div className="opacity-80">DOW: {dow.toLocaleString(undefined, { maximumFractionDigits: 2 })} pts</div>
            <div className="opacity-80">VIX: {vix.toLocaleString(undefined, { maximumFractionDigits: 2 })}</div>
            <div className="opacity-80">10Y: {y10.toLocaleString(undefined, { maximumFractionDigits: 2 })}%</div>
          </div>
          <div className="flex items-center gap-2 text-xs">
            {(["1d", "1w", "1m"] as RangeKey[]).map((k) => (
              <button
                key={k}
                onClick={() => setRange(k)}
                className={`rounded px-2 py-1 btn-sheen ${range === k ? "glass-card shadow-elevated" : "hover-lift"}`}
                style={{ color: "#0B2545" }}
              >
                {k.toUpperCase()}
              </button>
            ))}
          </div>
        </div>

        <div className="mt-2 grid grid-cols-5 gap-4 items-center">
          <div className="col-span-4 glass-card p-2 tilt-hover" style={{ width: "100%", height: 92 }}>
            <ResponsiveContainer>
              <BarChart data={spBars} margin={{ top: 2, bottom: 0, left: 0, right: 0 }}>
                <defs>
                  <linearGradient id="mkSpBar" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="10%" stopColor="#1E6AE1" stopOpacity={0.9} />
                    <stop offset="90%" stopColor="#1E6AE1" stopOpacity={0.2} />
                  </linearGradient>
                </defs>
                <XAxis dataKey="x" hide />
                <YAxis hide domain={spDomain as any} />
                <Tooltip formatter={(v: any) => [Number(v).toLocaleString(undefined, { maximumFractionDigits: 2 }), "Index"]} />
                <Bar dataKey="y" fill="url(#mkSpBar)" radius={[3,3,0,0]} barSize={5} isAnimationActive />
              </BarChart>
            </ResponsiveContainer>
          </div>
          <div className="col-span-1 text-right text-xs">
            <div style={{ color: trendColor }}>{trendArrow} {spDelta.delta.toFixed(2)}</div>
            <div className="opacity-80" style={{ color: "#374151" }}>{spDelta.pct.toFixed(2)}%</div>
            <div className="mt-2">
              <div className="h-1.5 bg-white/50 rounded relative overflow-hidden">
                <div className="absolute inset-y-0 left-0 bg-gradient-to-r from-[#1E6AE1] to-[#2AB3A6]" style={{ width: `${Math.max(0, Math.min(100, spPosPct))}%` }} />
              </div>
              <div className="mt-1 text-[10px]" style={{ color: "#6B7280" }}>Range pos (period)</div>
            </div>
          </div>
        </div>

        {loading && (
          <div className="mt-2 hero-progress"><span /></div>
        )}

        <div className="mt-3 grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div className="glass-card p-3 tilt-hover">
            <div className="text-xs" style={{ color: "#374151" }}>NASDAQ trend</div>
            <div style={{ width: "100%", height: 72 }}>
              <ResponsiveContainer>
                <BarChart data={ndqBars} margin={{ top: 2, bottom: 0, left: 0, right: 0 }}>
                  <defs>
                    <linearGradient id="mkNdqBar" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="10%" stopColor="#2AB3A6" stopOpacity={0.9} />
                      <stop offset="90%" stopColor="#2AB3A6" stopOpacity={0.2} />
                    </linearGradient>
                  </defs>
                  <XAxis dataKey="x" hide />
                  <YAxis hide domain={ndqDomain as any} />
                  <Tooltip formatter={(v: any) => [Number(v).toLocaleString(undefined, { maximumFractionDigits: 2 }), "Index"]} />
                  <Bar dataKey="y" fill="url(#mkNdqBar)" radius={[3,3,0,0]} barSize={5} isAnimationActive />
                </BarChart>
              </ResponsiveContainer>
            </div>
            <div className="mt-1 text-xs" style={{ color: ndqDelta.delta > 0 ? '#1D9A6C' : ndqDelta.delta < 0 ? '#E25555' : '#374151' }}>
              {(ndqDelta.delta > 0 ? '▲' : ndqDelta.delta < 0 ? '▼' : '')} {ndqDelta.delta.toFixed(2)} ({ndqDelta.pct.toFixed(2)}%)
            </div>
          </div>
          <div className="glass-card p-3 tilt-hover">
            <div className="text-xs" style={{ color: "#374151" }}>DOW trend</div>
            <div style={{ width: "100%", height: 72 }}>
              <ResponsiveContainer>
                <BarChart data={dowBars} margin={{ top: 2, bottom: 0, left: 0, right: 0 }}>
                  <defs>
                    <linearGradient id="mkDowBar" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="10%" stopColor="#9AD36A" stopOpacity={0.9} />
                      <stop offset="90%" stopColor="#9AD36A" stopOpacity={0.2} />
                    </linearGradient>
                  </defs>
                  <XAxis dataKey="x" hide />
                  <YAxis hide domain={dowDomain as any} />
                  <Tooltip formatter={(v: any) => [Number(v).toLocaleString(undefined, { maximumFractionDigits: 2 }), "Index"]} />
                  <Bar dataKey="y" fill="url(#mkDowBar)" radius={[3,3,0,0]} barSize={5} isAnimationActive />
                </BarChart>
              </ResponsiveContainer>
            </div>
            <div className="mt-1 text-xs" style={{ color: dowDelta.delta > 0 ? '#1D9A6C' : dowDelta.delta < 0 ? '#E25555' : '#374151' }}>
              {(dowDelta.delta > 0 ? '▲' : dowDelta.delta < 0 ? '▼' : '')} {dowDelta.delta.toFixed(2)} ({dowDelta.pct.toFixed(2)}%)
            </div>
          </div>
        </div>

        <div className="mt-3 grid grid-cols-2 sm:grid-cols-4 gap-3 text-xs">
          <div className="glass-card p-3 hover-lift">
            <div className="" style={{ color: "#374151" }}>VIX trend</div>
            <div style={{ width: "100%", height: 64 }}>
              <ResponsiveContainer>
                <BarChart data={vixBars} margin={{ top: 2, bottom: 0, left: 0, right: 0 }}>
                  <XAxis dataKey="x" hide />
                  <YAxis hide domain={paddedDomain(vixBars) as any} />
                  <Bar dataKey="y" fill="#0B2545" radius={[2,2,0,0]} barSize={5} isAnimationActive />
                </BarChart>
              </ResponsiveContainer>
            </div>
            <div className="mt-1" style={{ color: vixDelta > 0 ? '#E25555' : vixDelta < 0 ? '#1D9A6C' : '#374151' }}>
              {(vixDelta > 0 ? '▲' : vixDelta < 0 ? '▼' : '')} {vixDelta.toFixed(2)} ({vixPct.toFixed(2)}%)
            </div>
          </div>
          <div className="glass-card p-3 hover-lift">
            <div className="" style={{ color: "#374151" }}>10Y trend</div>
            <div style={{ width: "100%", height: 64 }}>
              <ResponsiveContainer>
                <BarChart data={y10Bars} margin={{ top: 2, bottom: 0, left: 0, right: 0 }}>
                  <XAxis dataKey="x" hide />
                  <YAxis hide domain={paddedDomain(y10Bars) as any} />
                  <Bar dataKey="y" fill="#7C3AED" radius={[2,2,0,0]} barSize={5} isAnimationActive />
                </BarChart>
              </ResponsiveContainer>
            </div>
            <div className="mt-1" style={{ color: y10Delta > 0 ? '#E25555' : y10Delta < 0 ? '#1D9A6C' : '#374151' }}>
              {(y10Delta > 0 ? '▲' : y10Delta < 0 ? '▼' : '')} {y10Delta.toFixed(3)}% ({y10Pct.toFixed(2)}%)
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}


