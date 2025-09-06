"use client";
import React, { useEffect, useMemo, useState } from "react";
import type { MarketLatestResponse, MarketHistoryResponse } from "../../types/api";
import { api } from "../../services/api";
import { LineChart, Line, ResponsiveContainer } from "recharts";
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
  const trendColor = spDelta.delta > 0 ? "#9AD36A" : spDelta.delta < 0 ? "#E25555" : "#FFFFFF";
  const trendArrow = spDelta.delta > 0 ? "▲" : spDelta.delta < 0 ? "▼" : "";

  const sp = useCountUp(m.sp500_close ?? 0);
  const ndq = useCountUp(m.nasdaq_close ?? 0);
  const dow = useCountUp(m.dow_close ?? 0);
  const vix = useCountUp(m.vix_close ?? 0);
  const y10 = useCountUp(m.treasury_10y ?? 0);

  return (
    <section className="animated-gradient">
      <div className="max-w-6xl mx-auto px-4 py-3 text-white text-shadow-soft" style={{ transition: "filter var(--dur-hover) var(--ease-enter)" }}>
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-4 text-sm">
            <div className="font-medium">Market Snapshot</div>
            <div className="opacity-95">S&amp;P 500: {sp.toLocaleString(undefined, { maximumFractionDigits: 2 })} pts</div>
            <div className="opacity-95">NASDAQ: {ndq.toLocaleString(undefined, { maximumFractionDigits: 2 })} pts</div>
            <div className="opacity-95">DOW: {dow.toLocaleString(undefined, { maximumFractionDigits: 2 })} pts</div>
            <div className="opacity-95">VIX: {vix.toLocaleString(undefined, { maximumFractionDigits: 2 })}</div>
            <div className="opacity-95">10Y: {y10.toLocaleString(undefined, { maximumFractionDigits: 2 })}%</div>
          </div>
          <div className="flex items-center gap-2 text-xs">
            {(["1d", "1w", "1m"] as RangeKey[]).map((k) => (
              <button
                key={k}
                onClick={() => setRange(k)}
                className={`rounded px-2 py-1 ${range === k ? "glass-card" : ""}`}
                style={{ color: "#FFFFFF" }}
              >
                {k.toUpperCase()}
              </button>
            ))}
          </div>
        </div>

        <div className="mt-2 grid grid-cols-5 gap-4 items-center">
          <div className="col-span-4" style={{ width: "100%", height: 42 }}>
            <ResponsiveContainer>
              <LineChart data={seriesSP} margin={{ top: 2, bottom: 0, left: 0, right: 0 }}>
                <Line type="monotone" dataKey="y" stroke="#FFFFFF" strokeOpacity={0.9} strokeWidth={1.5} dot={false} />
              </LineChart>
            </ResponsiveContainer>
          </div>
          <div className="col-span-1 text-right text-xs">
            <div style={{ color: trendColor }}>{trendArrow} {spDelta.delta.toFixed(2)}</div>
            <div className="opacity-80">{spDelta.pct.toFixed(2)}%</div>
          </div>
        </div>

        {loading && (
          <div className="mt-2 hero-progress"><span /></div>
        )}

        <div className="mt-3 grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div className="glass-card p-3">
            <div className="text-xs opacity-80">NASDAQ trend</div>
            <div style={{ width: "100%", height: 36 }}>
              <ResponsiveContainer>
                <LineChart data={seriesNDQ} margin={{ top: 2, bottom: 0, left: 0, right: 0 }}>
                  <Line type="monotone" dataKey="y" stroke="#FFFFFF" strokeOpacity={0.9} strokeWidth={1.2} dot={false} />
                </LineChart>
              </ResponsiveContainer>
            </div>
            <div className="mt-1 text-xs" style={{ color: ndqDelta.delta > 0 ? '#9AD36A' : ndqDelta.delta < 0 ? '#E25555' : '#FFFFFF' }}>
              {(ndqDelta.delta > 0 ? '▲' : ndqDelta.delta < 0 ? '▼' : '')} {ndqDelta.delta.toFixed(2)} ({ndqDelta.pct.toFixed(2)}%)
            </div>
          </div>
          <div className="glass-card p-3">
            <div className="text-xs opacity-80">DOW trend</div>
            <div style={{ width: "100%", height: 36 }}>
              <ResponsiveContainer>
                <LineChart data={seriesDOW} margin={{ top: 2, bottom: 0, left: 0, right: 0 }}>
                  <Line type="monotone" dataKey="y" stroke="#FFFFFF" strokeOpacity={0.9} strokeWidth={1.2} dot={false} />
                </LineChart>
              </ResponsiveContainer>
            </div>
            <div className="mt-1 text-xs" style={{ color: dowDelta.delta > 0 ? '#9AD36A' : dowDelta.delta < 0 ? '#E25555' : '#FFFFFF' }}>
              {(dowDelta.delta > 0 ? '▲' : dowDelta.delta < 0 ? '▼' : '')} {dowDelta.delta.toFixed(2)} ({dowDelta.pct.toFixed(2)}%)
            </div>
          </div>
        </div>

        <div className="mt-3 grid grid-cols-2 sm:grid-cols-4 gap-3 text-xs">
          <div className="glass-card p-3">
            <div className="opacity-80">VIX Δ</div>
            <div style={{ color: vixDelta > 0 ? '#E25555' : vixDelta < 0 ? '#9AD36A' : '#FFFFFF' }}>
              {(vixDelta > 0 ? '▲' : vixDelta < 0 ? '▼' : '')} {vixDelta.toFixed(2)} ({vixPct.toFixed(2)}%)
            </div>
          </div>
          <div className="glass-card p-3">
            <div className="opacity-80">10Y Δ</div>
            <div style={{ color: y10Delta > 0 ? '#E25555' : y10Delta < 0 ? '#9AD36A' : '#FFFFFF' }}>
              {(y10Delta > 0 ? '▲' : y10Delta < 0 ? '▼' : '')} {y10Delta.toFixed(3)}% ({y10Pct.toFixed(2)}%)
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}


