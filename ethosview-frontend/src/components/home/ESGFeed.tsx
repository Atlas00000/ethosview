"use client";
import React from "react";
import type { ESGScoresListResponse } from "../../types/api";
import { CompanyQuickView } from "./CompanyQuickView";
import { useCountUp } from "./useCountUp";

function CountUpValue({ value }: { value: number }) {
  const count = useCountUp(typeof value === 'number' ? value : 0);
  return <>{Number.isFinite(count as any) ? (count as number).toFixed(2) : '—'}</>;
}

export function ESGFeed({ list }: { list: ESGScoresListResponse }) {
  const [items, setItems] = React.useState(list?.scores || []);
  const [offset, setOffset] = React.useState((list?.pagination?.limit ?? 20));
  const [loading, setLoading] = React.useState(false);
  const [quickId, setQuickId] = React.useState<number | null>(null);
  const listRef = React.useRef<HTMLDivElement | null>(null);
  if (!items.length) return null;
  const best = [...items].sort((a, b) => (b.overall_score ?? 0) - (a.overall_score ?? 0)).slice(0, 5);
  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-3">
        <h2 className="text-xl font-semibold text-gradient">Latest ESG scores</h2>
        <span className="text-xs" style={{ color: "#374151" }}>Showing {items.length} items</span>
      </div>
      <div ref={listRef} className="grid grid-cols-1 lg:grid-cols-3 gap-3">
        <div className="glass-card divide-y animate-fade-in-up lg:col-span-2">
          {items.map((s) => (
            <div key={s.id} className="p-3 flex items-center justify-between hover-lift tilt-hover">
              <div className="truncate pr-3">
                <div className="text-sm font-medium" style={{ color: "#0B2545" }}>{s.company_name || s.company_id}</div>
                <div className="text-xs" style={{ color: "#374151" }}>{s.score_date ? new Date(s.score_date).toISOString().slice(0,10) : 'N/A'} {s.company_symbol ? `• ${s.company_symbol}` : ''}</div>
              </div>
              <div className="flex items-center gap-2">
                <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}><CountUpValue value={typeof s.overall_score === 'number' ? s.overall_score : 0} /></div>
                <button className="glass-card px-2 py-1 btn-sheen text-xs" onClick={() => setQuickId(s.company_id)}>View</button>
              </div>
            </div>
          ))}
        </div>
        <div className="glass-card p-3 animate-fade-in-up">
          <div className="text-sm font-medium mb-2" style={{ color: "#0B2545" }}>Top latest</div>
          <div className="space-y-2">
            {best.map((s) => (
              <div key={`best-${s.id}`} className="glass-card p-2 hover-lift">
                <div className="flex items-center justify-between text-sm">
                  <div className="truncate pr-2" style={{ color: "#0B2545" }}>{s.company_name || s.company_id}</div>
                  <div className="font-semibold" style={{ color: "#1D9A6C" }}>{(s.overall_score ?? 0).toFixed(2)}</div>
                </div>
                <div className="mt-1 h-1.5 bg-white/30 rounded">
                  <div className="h-1.5 rounded" style={{ width: `${Math.max(0, Math.min(100, s.overall_score ?? 0))}%`, background: '#1D9A6C', transition: 'width 320ms var(--ease-enter)' }} />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
      <div className="mt-4 flex justify-center">
        <button
          disabled={loading}
          className="btn-primary btn-sheen px-4 py-2 rounded hover-lift"
          onClick={async () => {
            try {
              setLoading(true);
              const next = await (await import("../../services/api")).api.esgScores(20, 0, offset);
              const scores = next?.scores || [];
              setItems((prev) => {
                const seen = new Set(prev.map(x => x.id));
                const merged = [...prev];
                for (const s of scores) if (!seen.has(s.id)) merged.push(s);
                return merged;
              });
              setOffset((o) => o + (next?.pagination?.limit ?? 20));
              // Prefetch next page when scrolled near bottom
              requestIdleCallback(() => {
                const el = listRef.current;
                if (!el) return;
                const { scrollHeight, clientHeight } = el;
                if (scrollHeight <= clientHeight * 1.6) {
                  (async () => {
                    const more = await (await import("../../services/api")).api.esgScores(20, 0, offset + (next?.pagination?.limit ?? 20)).catch(() => null);
                    if (!more?.scores?.length) return;
                    setItems((prev) => {
                      const seen2 = new Set(prev.map(x => x.id));
                      const merged2 = [...prev];
                      for (const s2 of more.scores) if (!seen2.has(s2.id)) merged2.push(s2);
                      return merged2;
                    });
                    setOffset((o) => o + (more?.pagination?.limit ?? 20));
                  })();
                }
              });
            } finally {
              setLoading(false);
            }
          }}
        >
          {loading ? "Loading…" : "Load more"}
        </button>
      </div>
      {quickId !== null && <CompanyQuickView companyId={quickId} onClose={() => setQuickId(null)} />}
    </section>
  );
}


