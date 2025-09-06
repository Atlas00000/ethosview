"use client";
import React from "react";
import type { TopPEResponse } from "../../types/api";
import { useCountUp } from "./useCountUp";

export function PELeaders({ top }: { top: TopPEResponse }) {
  const items = top.top_performers;
  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <h2 className="text-xl font-semibold mb-3 text-gradient">Top P/E (lower is better)</h2>
      <div className="text-xs mb-2" style={{ color: "#374151" }}>P/E reflects price relative to earnings; lower can indicate value.</div>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        {items.slice(0, 10).map((it) => {
          const v = typeof it.value === 'number' ? it.value : 0;
          const count = useCountUp(v);
          return (
            <div key={it.company_id} className="glass-card p-3 hover-lift animate-fade-in-up">
              <div className="flex items-center justify-between">
                <div>
                  <div className="text-sm font-medium" style={{ color: "#0B2545" }}>{it.company_name}</div>
                  <div className="text-xs" style={{ color: "#374151" }}>Rank #{it.rank}</div>
                </div>
                <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{Number.isFinite(count as any) ? (count as number).toFixed(2) : '—'}</div>
              </div>
              <div className="mt-2 h-1.5 bg-white/30 rounded">
                <div className="h-1.5 rounded" style={{ width: `${Math.min(100, (100 / Math.max(1, v))) }%`, background: '#1E6AE1' }} />
              </div>
            </div>
          );
        })}
      </div>
    </section>
  );
}


