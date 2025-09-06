"use client";
import React, { useMemo, useState } from "react";
import type { ESGTrendsResponse } from "../../types/api";
import { LineChart, Line, ResponsiveContainer, Tooltip, XAxis, YAxis, AreaChart, Area } from "recharts";
import { useCountUp } from "./useCountUp";

export function ESGTrendMini({ trends }: { trends: ESGTrendsResponse }) {
  const rows = (trends?.trends || []).slice().reverse();
  const data = rows.map((t) => ({ x: t.date, y: t.esg_score, e: t.e_score, s: t.s_score, g: t.g_score }));
  const last = data.length ? data[data.length - 1] : null;
  const cESG = useCountUp(typeof last?.y === 'number' ? last!.y : 0);
  const [mode, setMode] = useState<'overall' | 'components'>('overall');
  const avg = useMemo(() => {
    if (!data.length) return null;
    const sum = data.reduce((acc, d) => acc + (d.y || 0), 0);
    return sum / data.length;
  }, [data]);
  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <h2 className="text-xl font-semibold mb-3 text-gradient">ESG trend</h2>
      <div className="glass-card p-3 animate-fade-in-up" style={{ width: "100%", height: 160 }}>
        <div className="flex items-center justify-between mb-1">
          <div className="text-xs" style={{ color: "#374151" }}>Latest ESG</div>
          <div className="flex items-center gap-2 text-xs">
            <button className={`px-2 py-0.5 rounded-full ${mode==='overall'?'bg-white/60 shadow-elevated':'bg-white/30 hover:bg-white/50'} tilt-hover`} onClick={()=>setMode('overall')}>Overall</button>
            <button className={`px-2 py-0.5 rounded-full ${mode==='components'?'bg:white/60 shadow-elevated':'bg-white/30 hover:bg-white/50'} tilt-hover`} onClick={()=>setMode('components')}>E/S/G</button>
          </div>
        </div>
        <div className="grid grid-cols-3 gap-2 mb-2">
          <div className="glass-card p-2 text-center">
            <div className="text-[10px]" style={{ color: "#374151" }}>Current</div>
            <div className="text-base font-semibold" style={{ color: "#1D9A6C" }}>{(typeof cESG==='number' && isFinite(cESG)) ? cESG.toFixed(1) : '—'}</div>
          </div>
          <div className="glass-card p-2 text-center">
            <div className="text-[10px]" style={{ color: "#374151" }}>Average</div>
            <div className="text-base font-semibold" style={{ color: "#0B2545" }}>{(typeof avg==='number' && isFinite(avg)) ? avg.toFixed(1) : '—'}</div>
          </div>
          <div className="glass-card p-2 text-center">
            <div className="text-[10px]" style={{ color: "#374151" }}>Period</div>
            <div className="text-base font-semibold" style={{ color: "#0B2545" }}>{data.length} pts</div>
          </div>
        </div>
        <ResponsiveContainer>
          {mode === 'overall' ? (
            <AreaChart data={data} margin={{ top: 4, bottom: 0, left: 0, right: 0 }}>
              <XAxis dataKey="x" hide />
              <YAxis domain={[0, 100]} hide />
              <Tooltip formatter={(v: any)=> (typeof v==='number'? v.toFixed(2): v)} labelFormatter={(l)=>`Date: ${l}`} />
              <Area type="monotone" dataKey="y" stroke="#1D9A6C" fill="#1D9A6C22" strokeWidth={1.5} />
            </AreaChart>
          ) : (
            <LineChart data={data} margin={{ top: 4, bottom: 0, left: 0, right: 0 }}>
              <XAxis dataKey="x" hide />
              <YAxis domain={[0, 100]} hide />
              <Tooltip formatter={(v: any)=> (typeof v==='number'? v.toFixed(2): v)} labelFormatter={(l)=>`Date: ${l}`} />
              <Line type="monotone" dataKey="e" stroke="#2AB3A6" strokeWidth={1.2} dot={false} />
              <Line type="monotone" dataKey="s" stroke="#1E6AE1" strokeWidth={1.2} dot={false} />
              <Line type="monotone" dataKey="g" stroke="#9AD36A" strokeWidth={1.2} dot={false} />
            </LineChart>
          )}
        </ResponsiveContainer>
      </div>
    </section>
  );
}


