"use client";
import React, { useMemo, useState } from "react";
import { api } from "../../services/api";
import type { CompanyResponse, CompanyFinancialSummaryResponse, LatestESGResponse, CompanyPricesResponse } from "../../types/api";
import { CompanySparkline } from "./MarketSparkline";

type LookupCard = {
  company: CompanyResponse;
  summary?: CompanyFinancialSummaryResponse | null;
  esg?: LatestESGResponse | null;
  loading: boolean;
  error?: string | null;
};

export function SymbolLookup() {
  const [symbols, setSymbols] = useState("");
  const [cards, setCards] = useState<LookupCard[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSearch(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setCards([]);
    const list = symbols
      .split(/[,\s]+/)
      .map((s) => s.trim().toUpperCase())
      .filter(Boolean)
      .slice(0, 8);
    if (!list.length) return;
    setLoading(true);
    try {
      // Resolve companies by symbol (hugging backend)
      const companies = (await Promise.all(
        list.map((sym) => api.companyBySymbol(sym).catch(() => null))
      )) as (CompanyResponse | null)[];
      const valid = companies.filter(Boolean) as CompanyResponse[];
      if (!valid.length) {
        setError("No companies found");
        return;
      }
      // seed cards
      setCards(valid.map((c) => ({ company: c, loading: true })));
      // Fetch details in parallel per company
      const rows = await Promise.all(
        valid.map(async (c) => {
          try {
            const [summary, esg, prices] = await Promise.all([
              api.companyFinancialSummary(c.id).catch(() => null),
              api.latestESG(c.id).catch(() => null),
              api.stockPrices(c.id, 30).catch(() => null),
            ]);
            return { id: c.id, summary: summary as CompanyFinancialSummaryResponse | null, esg: esg as LatestESGResponse | null, prices: prices as CompanyPricesResponse | null };
          } catch {
            return { id: c.id, summary: null, esg: null };
          }
        })
      );
      setCards((prev) =>
        prev.map((card) => {
          const r = rows.find((x) => x.id === card.company.id);
          return { ...card, summary: r?.summary ?? null, esg: r?.esg ?? null, loading: false, prices: (r as any)?.prices ?? null };
        })
      );
    } catch (err: any) {
      setError("Lookup failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="max-w-6xl mx-auto px-4 py-8 animate-fade-in-up">
      <div className="glass-card p-4 shadow-elevated">
        <form onSubmit={onSearch} className="flex flex-wrap gap-2 items-center">
          <input
            className="border rounded px-3 py-2 w-64 focus:outline-none focus:ring-2 hover-lift"
            placeholder="Lookup by symbol(s), e.g., AAPL, MSFT, TSLA"
            value={symbols}
            onChange={(e) => setSymbols(e.target.value)}
          />
          <button className="btn-primary rounded px-3 py-2 btn-sheen" disabled={loading}>
            {loading ? "Searching..." : "Search"}
          </button>
          <div className="text-xs" style={{ color: "#374151" }}>
            Tip: separate multiple symbols with commas or spaces
          </div>
        </form>
        {error && <div className="text-sm text-red-600 mt-2">{error}</div>}
      </div>

      {cards.length > 0 && (
        <div className="mt-4 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
          {cards.map((card) => {
            const s = card.summary?.summary;
            const ind = card.summary?.indicators;
            const price = s?.current_price;
            const changePct = s?.price_change_percent;
            const pe = ind?.pe_ratio;
            const mcap = ind?.market_cap;
            const esg = card.esg?.overall_score;
            return (
              <div key={card.company.id} className="group glass-card p-4 hover-lift tilt-hover animate-fade-in-up relative overflow-hidden">
                <div className="absolute inset-0 pointer-events-none button-sheen opacity-0 group-hover:opacity-100" />
                <div className="flex items-center justify-between">
                  <div>
                    <div className="text-sm font-medium" style={{ color: "#0B2545" }}>{card.company.name} ({card.company.symbol})</div>
                    <div className="text-xs" style={{ color: "#374151" }}>{card.company.sector} • {card.company.country}</div>
                  </div>
                  <span className="px-2 py-0.5 rounded-full text-[10px] bg-white/60">ID #{card.company.id}</span>
                </div>
                <div className="mt-3 grid grid-cols-3 gap-2">
                  <div className="glass-card p-2 text-center">
                    <div className="text-[10px]" style={{ color: "#374151" }}>Price</div>
                    <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{card.loading ? <span className="inline-block w-10 h-4 skeleton rounded" /> : (typeof price === 'number' ? price.toFixed(2) : '—')}</div>
                  </div>
                  <div className="glass-card p-2 text-center">
                    <div className="text-[10px]" style={{ color: "#374151" }}>Change</div>
                    <div className={`text-sm font-semibold ${typeof changePct === 'number' ? (changePct >= 0 ? 'text-green-600' : 'text-red-600') : ''}`}>{card.loading ? <span className="inline-block w-10 h-4 skeleton rounded" /> : (typeof changePct === 'number' ? `${changePct >= 0 ? '+' : ''}${changePct.toFixed(2)}%` : '—')}</div>
                  </div>
                  <div className="glass-card p-2 text-center">
                    <div className="text-[10px]" style={{ color: "#374151" }}>P/E</div>
                    <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{card.loading ? <span className="inline-block w-8 h-4 skeleton rounded" /> : (typeof pe === 'number' ? pe.toFixed(1) : '—')}</div>
                  </div>
                </div>
                <div className="mt-2 grid grid-cols-3 gap-2">
                  <div className="glass-card p-2 text-center">
                    <div className="text-[10px]" style={{ color: "#374151" }}>Market Cap</div>
                    <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{card.loading ? <span className="inline-block w-16 h-4 skeleton rounded" /> : (typeof mcap === 'number' ? formatCompact(mcap) : '—')}</div>
                  </div>
                  <div className="glass-card p-2 text-center">
                    <div className="text-[10px]" style={{ color: "#374151" }}>ESG</div>
                    <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{card.loading ? <span className="inline-block w-8 h-4 skeleton rounded" /> : (typeof esg === 'number' ? esg.toFixed(0) : '—')}</div>
                  </div>
                  <div className="glass-card p-2 text-center">
                    <div className="text-[10px]" style={{ color: "#374151" }}>Industry</div>
                    <div className="text-xs font-medium" style={{ color: "#0B2545" }}>{card.company.industry || '—'}</div>
                  </div>
                </div>
                <div className="mt-2">
                  {card.prices && <CompanySparkline series={card.prices} />}
                </div>
              </div>
            );
          })}
        </div>
      )}
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


