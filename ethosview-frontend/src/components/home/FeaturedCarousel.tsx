"use client";
import React, { useEffect, useMemo, useRef, useState } from "react";
import type { PerformanceMetric, CompanyFinancialSummaryResponse } from "../../types/api";
import { api } from "../../services/api";
import { useCountUp } from "./useCountUp";
import { CompanySparkline } from "./MarketSparkline";
import { CompanyQuickView } from "./CompanyQuickView";

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

  const [details, setDetails] = useState<Record<number, { price?: number; esg?: number; pe?: number; mcap?: number; changePct?: number; prices?: any }>>({});
  const [activeIndex, setActiveIndex] = useState(0);
  const scrollerRef = useRef<HTMLDivElement | null>(null);
  const [quickId, setQuickId] = useState<number | null>(null);
  const observerRef = useRef<IntersectionObserver | null>(null);

  // Lazy-load details as cards enter viewport (reduce burst, enable more items)
  useEffect(() => {
    if (observerRef.current) observerRef.current.disconnect();
    observerRef.current = new IntersectionObserver((entries) => {
      entries.forEach(async (entry) => {
        if (!entry.isIntersecting) return;
        const el = entry.target as HTMLElement;
        const idAttr = el.getAttribute('data-id');
        if (!idAttr) return;
        const id = Number(idAttr);
        if (!id || details[id]) return;
        try {
          const [p, e, s, hist] = await Promise.all([
            api.latestPrice(id).catch(() => null),
            api.latestESG(id).catch(() => null),
            api.companyFinancialSummary(id).catch(() => null),
            api.stockPrices(id, 30).catch(() => null),
          ]);
          const sum = s as CompanyFinancialSummaryResponse | null;
          setDetails((prev) => ({
            ...prev,
            [id]: {
              price: (p as any)?.price?.close_price,
              esg: (e as any)?.overall_score,
              pe: sum?.indicators?.pe_ratio,
              mcap: sum?.indicators?.market_cap,
              changePct: sum?.summary?.price_change_percent,
              prices: hist,
            },
          }));
        } catch {}
      });
    }, { root: scrollerRef.current, rootMargin: '0px 200px', threshold: 0.2 });
    const root = scrollerRef.current;
    const cards = root?.querySelectorAll('[data-id]');
    cards?.forEach((c) => observerRef.current?.observe(c));
    return () => observerRef.current?.disconnect();
  }, [items.length]);

  // Auto-advance carousel subtly
  useEffect(() => {
    const el = scrollerRef.current;
    if (!el) return;
    const timer = setInterval(() => {
      const next = (activeIndex + 1) % items.length;
      setActiveIndex(next);
      el.scrollTo({ left: next * 260, behavior: "smooth" });
    }, 5000);
    return () => clearInterval(timer);
  }, [activeIndex, items.length]);

  return (
    <section className="max-w-6xl mx-auto px-4 py-12 relative overflow-hidden">
      <div className="blob" style={{ bottom: -80, left: -60, width: 200, height: 200, background: "rgba(30,106,225,0.12)", borderRadius: 9999 }} />
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold" style={{ color: "#0B2545" }}>Featured companies</h2>
        <div className="flex items-center gap-2">
          <button aria-label="Prev" className="glass-card px-3 py-1 btn-sheen" onClick={() => {
            const prev = (activeIndex - 1 + items.length) % items.length;
            setActiveIndex(prev);
            scrollerRef.current?.scrollTo({ left: prev * 260, behavior: "smooth" });
          }}>‹</button>
          <button aria-label="Next" className="glass-card px-3 py-1 btn-sheen" onClick={() => {
            const next = (activeIndex + 1) % items.length;
            setActiveIndex(next);
            scrollerRef.current?.scrollTo({ left: next * 260, behavior: "smooth" });
          }}>›</button>
        </div>
      </div>
      <div ref={scrollerRef} className="flex gap-4 overflow-x-auto pb-2 edge-fade snap-x snap-mandatory">
        {items.map((it, idx) => {
          const d = details[it.company_id] || {};
          const p = typeof d.price === 'number' ? d.price : undefined;
          const esg = typeof d.esg === 'number' ? d.esg : undefined;
          const pe = typeof d.pe === 'number' ? d.pe : undefined;
          const mcap = typeof d.mcap === 'number' ? d.mcap : undefined;
          const changePct = typeof d.changePct === 'number' ? d.changePct : undefined;
          const pCount = useCountUp(p ?? 0);
          const esgCount = useCountUp(esg ?? 0);
          const strength = typeof esg === 'number' ? (esg >= 75 ? 'Strong' : esg >= 50 ? 'Moderate' : 'Developing') : '—';
          return (
            <div key={it.company_id} data-id={it.company_id} className="group min-w-[240px] shrink-0 glass-card p-4 hover-lift tilt-hover snap-start animate-fade-in-up relative overflow-hidden shadow-elevated"
                 style={{ transition: "transform 240ms var(--ease-enter), box-shadow 240ms var(--ease-enter)" }}>
              <div className="absolute inset-0 pointer-events-none button-sheen opacity-0 group-hover:opacity-100" />
              <div className="text-sm" style={{ color: "#6B7280" }}>#{it.company_id}</div>
              <div className="text-lg font-medium" style={{ color: "#0B2545" }}>{it.company_name}</div>
              <div className="mt-2 grid grid-cols-2 gap-2">
                <div className="glass-card p-2 text-center float-slow">
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
              <div className="mt-2 grid grid-cols-2 gap-2">
                <div className="glass-card p-2 text-center">
                  <div className="text-[10px]" style={{ color: "#374151" }}>P/E</div>
                  <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{pe ? pe.toFixed(1) : '—'}</div>
                </div>
                <div className="glass-card p-2 text-center">
                  <div className="text-[10px]" style={{ color: "#374151" }}>Mkt Cap</div>
                  <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{mcap ? formatCompact(mcap) : '—'}</div>
                </div>
              </div>
              <div className="mt-2 flex items-center justify-between">
                <div className={`text-xs font-medium ${changePct && changePct >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                  {typeof changePct === 'number' ? `${changePct >= 0 ? '+' : ''}${changePct.toFixed(2)}%` : '—'}
                </div>
                <button className="btn-primary btn-sheen px-3 py-1 rounded text-xs pulse-outline" onClick={() => setQuickId(it.company_id)}>Details</button>
              </div>
              <div className="mt-2">
                <CompanySparkline series={d.prices ?? null} />
              </div>
            </div>
          );
        })}
      </div>
      {quickId !== null && <CompanyQuickView companyId={quickId} onClose={() => setQuickId(null)} />}
    </section>
  );
}


function formatCompact(n: number) {
  if (n >= 1e12) return (n / 1e12).toFixed(2) + 'T';
  if (n >= 1e9) return (n / 1e9).toFixed(2) + 'B';
  if (n >= 1e6) return (n / 1e6).toFixed(2) + 'M';
  if (n >= 1e3) return (n / 1e3).toFixed(2) + 'K';
  return n.toFixed(0);
}

