"use client";
import React, { useEffect, useMemo, useState } from "react";
import type { PerformanceMetric } from "../../types/api";
import { api } from "../../services/api";
import { useCountUp } from "./useCountUp";

type FeaturedItem = {
  company_id: number;
  company_name: string;
  price?: number;
  esg?: number;
};

export function FeaturedCarousel({ top }: { top: PerformanceMetric[] }) {
  const items: FeaturedItem[] = top.map((t) => ({
    company_id: t.company_id,
    company_name: t.company_name,
  }));

  const [details, setDetails] = useState<Record<number, { price?: number; esg?: number }>>({});

  useEffect(() => {
    const ids = items.slice(0, 6).map((i) => i.company_id);
    if (!ids.length) return;
    Promise.all(ids.map(async (id) => {
      try {
        const [p, e] = await Promise.all([
          api.latestPrice(id).catch(() => null),
          api.latestESG(id).catch(() => null),
        ]);
        return { id, price: p?.price?.close_price, esg: (e as any)?.overall_score };
      } catch {
        return { id } as any;
      }
    })).then((rows) => {
      setDetails((prev) => {
        const next = { ...prev } as Record<number, { price?: number; esg?: number }>;
        for (const r of rows) next[r.id] = { price: r.price, esg: r.esg };
        return next;
      });
    });
  }, [top]);

  return (
    <section className="max-w-6xl mx-auto px-4 py-12 relative overflow-hidden">
      <div className="blob" style={{ bottom: -80, left: -60, width: 200, height: 200, background: "rgba(30,106,225,0.12)", borderRadius: 9999 }} />
      <h2 className="text-xl font-semibold mb-4" style={{ color: "#0B2545" }}>Featured companies</h2>
      <div className="flex gap-4 overflow-x-auto pb-2 edge-fade snap-x snap-mandatory">
        {items.map((it) => {
          const d = details[it.company_id] || {};
          const p = typeof d.price === 'number' ? d.price : undefined;
          const esg = typeof d.esg === 'number' ? d.esg : undefined;
          const pCount = useCountUp(p ?? 0);
          const esgCount = useCountUp(esg ?? 0);
          const strength = typeof esg === 'number' ? (esg >= 75 ? 'Strong' : esg >= 50 ? 'Moderate' : 'Developing') : '—';
          return (
            <div key={it.company_id} className="group min-w-[240px] shrink-0 glass-card p-4 hover-lift tilt-hover snap-start animate-fade-in-up relative overflow-hidden">
              <div className="absolute inset-0 pointer-events-none button-sheen opacity-0 group-hover:opacity-100" />
              <div className="text-sm" style={{ color: "#6B7280" }}>#{it.company_id}</div>
              <div className="text-lg font-medium" style={{ color: "#0B2545" }}>{it.company_name}</div>
              <div className="mt-2 grid grid-cols-2 gap-2">
                <div className="glass-card p-2 text-center">
                  <div className="text-[10px]" style={{ color: "#374151" }}>Price</div>
                  <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{p !== undefined ? pCount.toFixed(2) : '—'}</div>
                </div>
                <div className="glass-card p-2 text-center">
                  <div className="text-[10px]" style={{ color: "#374151" }}>ESG</div>
                  <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{esg !== undefined ? esgCount.toFixed(2) : '—'}</div>
                </div>
              </div>
              <div className="mt-2 flex items-center gap-2">
                <span className="px-2 py-0.5 rounded-full text-[10px] bg-white/50">{strength}</span>
                {typeof esg === 'number' && (
                  <div className="flex-1 h-1.5 bg-white/30 rounded">
                    <div className="h-1.5 rounded" style={{ width: `${Math.max(0, Math.min(100, esg))}%`, background: '#1D9A6C' }} />
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </section>
  );
}


