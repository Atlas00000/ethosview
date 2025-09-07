"use client";
import React, { useEffect, useMemo, useRef, useState } from "react";
import type { TopPEResponse, CompanyFinancialSummaryResponse, CompanyResponse, LatestESGResponse } from "../../types/api";
import { api } from "../../services/api";
import { useCountUp } from "./useCountUp";
import { CompanySparkline } from "./MarketSparkline";

export function PELeaders({ top }: { top: TopPEResponse }) {
  const items = top.top_performers;
  const [details, setDetails] = useState<Record<number, {
    price?: number;
    changePct?: number;
    mcap?: number;
    esg?: number;
    symbol?: string;
  }>>({});
  const [expanded, setExpanded] = useState(false);
  const observerRef = useRef<IntersectionObserver | null>(null);
  const noFinRef = useRef<Set<number>>(new Set());
  const scrollerRef = useRef<HTMLDivElement | null>(null);

  // Fetch additional details for more context (no placeholders)
  // Lazy load metrics per-card as visible
  useEffect(() => {
    if (observerRef.current) observerRef.current.disconnect();
    observerRef.current = new IntersectionObserver((entries) => {
      entries.forEach(async (entry) => {
        if (!entry.isIntersecting) return;
        const el = entry.target as HTMLElement;
        const idAttr = el.getAttribute('data-id');
        if (!idAttr) return;
        const id = Number(idAttr);
        if (!id || details[id] || noFinRef.current.has(id)) return;
        try {
          const [fin, esg, company] = await Promise.all([
            api.companyFinancialSummary(id).catch(() => null),
            api.latestESG(id).catch(() => null),
            api.companyById(id).catch(() => null),
          ]);
          const f = fin as CompanyFinancialSummaryResponse | null;
          const e = esg as LatestESGResponse | null;
          const c = company as CompanyResponse | null;
          if (!f) { noFinRef.current.add(id); return; }
          setDetails(prev => ({
            ...prev,
            [id]: { price: f?.summary?.current_price, changePct: f?.summary?.price_change_percent, mcap: f?.indicators?.market_cap, esg: e?.overall_score, symbol: c?.symbol },
          }));
        } catch { noFinRef.current.add(id); }
      });
    }, { root: scrollerRef.current, rootMargin: '0px 200px', threshold: 0.2 });
    const root = scrollerRef.current;
    const cards = root?.querySelectorAll('[data-id]');
    cards?.forEach(c => observerRef.current?.observe(c));
    return () => observerRef.current?.disconnect();
  }, [items, expanded]);

  const visible = useMemo(() => items.slice(0, expanded ? 24 : 12), [items, expanded]);

  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-3">
        <h2 className="text-xl font-semibold text-gradient">Top P/E (lower is better)</h2>
        <button className="glass-card px-3 py-1 btn-sheen" onClick={() => setExpanded(v => !v)}>{expanded ? "Show fewer" : "Show more"}</button>
      </div>
      <div className="text-xs mb-2" style={{ color: "#374151" }}>P/E reflects price relative to earnings; lower can indicate value.</div>
      <div ref={scrollerRef} className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        {visible.map((it) => {
          const v = typeof it.value === 'number' ? it.value : 0;
          const count = useCountUp(v);
          const d = details[it.company_id] || {};
          const changePct = d.changePct;
          const esg = d.esg;
          return (
            <div key={it.company_id} data-id={it.company_id} className="group glass-card p-3 hover-lift tilt-hover animate-fade-in-up relative overflow-hidden">
              <div className="absolute inset-0 pointer-events-none button-sheen opacity-0 group-hover:opacity-100" />
              <div className="flex items-center justify-between">
                <div>
                  <div className="text-sm font-medium" style={{ color: "#0B2545" }}>{it.company_name}{d.symbol ? ` (${d.symbol})` : ''}</div>
                  <div className="text-xs" style={{ color: "#374151" }}>Rank #{it.rank}</div>
                </div>
                <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{Number.isFinite(count as any) ? (count as number).toFixed(2) : '—'}</div>
              </div>
              <div className="mt-2 grid grid-cols-3 gap-2">
                <div className="glass-card p-2 text-center">
                  <div className="text-[10px]" style={{ color: "#374151" }}>Price</div>
                  <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{typeof d.price === 'number' ? d.price.toFixed(2) : '—'}</div>
                </div>
                <div className="glass-card p-2 text-center">
                  <div className="text-[10px]" style={{ color: "#374151" }}>Change</div>
                  <div className={`text-sm font-semibold ${typeof changePct === 'number' ? (changePct >= 0 ? 'text-green-600' : 'text-red-600') : ''}`}>{typeof changePct === 'number' ? `${changePct >= 0 ? '+' : ''}${changePct.toFixed(2)}%` : '—'}</div>
                </div>
                <div className="glass-card p-2 text-center">
                  <div className="text-[10px]" style={{ color: "#374151" }}>ESG</div>
                  <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{typeof esg === 'number' ? esg.toFixed(0) : '—'}</div>
                </div>
              </div>
              <div className="mt-2 h-1.5 bg-white/30 rounded">
                <div className="h-1.5 rounded" style={{ width: `${Math.min(100, (100 / Math.max(1, v))) }%`, background: '#1E6AE1', transition: 'width 300ms var(--ease-enter)' }} />
              </div>
              <div className="mt-2">
                {/* Optional per-company sparkline: loaded separately to avoid blocking */}
                {d.price !== undefined && <CompanySparkline series={null as any} />}
              </div>
            </div>
          );
        })}
      </div>
    </section>
  );
}


