"use client";
import React, { useEffect, useState } from "react";
import { api } from "../../services/api";
import type {
  CompanyResponse,
  LatestESGResponse,
  FinancialIndicatorsResponse,
  LatestPriceResponse,
  CompanyPricesResponse,
  CompanyFinancialSummaryResponse,
} from "../../types/api";
import { CompanySparkline } from "./MarketSparkline";

type Props = {
  companyId: number;
  onClose: () => void;
};

export function CompanyQuickView({ companyId, onClose }: Props) {
  const [company, setCompany] = useState<CompanyResponse | null>(null);
  const [esg, setEsg] = useState<LatestESGResponse | null>(null);
  const [ind, setInd] = useState<FinancialIndicatorsResponse | null>(null);
  const [price, setPrice] = useState<LatestPriceResponse | null>(null);
  const [series, setSeries] = useState<CompanyPricesResponse | null>(null);
  const [summary, setSummary] = useState<CompanyFinancialSummaryResponse | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let alive = true;
    setLoading(true);
    Promise.all([
      api.companyById(companyId).catch(() => null),
      api.latestESG(companyId).catch(() => null),
      api.financialIndicators(companyId).catch(() => null),
      api.latestPrice(companyId).catch(() => null),
      api.stockPrices(companyId, 30).catch(() => null),
      api.companyFinancialSummary(companyId).catch(() => null),
    ]).then(([c, e, i, p, s, sum]) => {
      if (!alive) return;
      setCompany(c as CompanyResponse | null);
      setEsg(e as LatestESGResponse | null);
      setInd(i as FinancialIndicatorsResponse | null);
      setPrice(p as LatestPriceResponse | null);
      setSeries(s as CompanyPricesResponse | null);
      setSummary(sum as CompanyFinancialSummaryResponse | null);
    }).finally(() => {
      if (alive) setLoading(false);
    });
    return () => {
      alive = false;
    };
  }, [companyId]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center" style={{ background: "rgba(11,37,69,0.35)", backdropFilter: "blur(6px)" }}>
      <div className="glass-card shadow-elevated animate-fade-in-up" style={{ width: "min(900px, 92vw)" }}>
        <div className="p-3 flex items-center justify-between border-b">
          <div className="text-sm font-medium" style={{ color: "#0B2545" }}>{company ? `${company.name} (${company.symbol})` : `Company #${companyId}`}</div>
          <button className="glass-card px-2 py-1 btn-sheen" onClick={onClose}>Close</button>
        </div>
        <div className="p-4 grid grid-cols-1 lg:grid-cols-3 gap-3">
          <div className="lg:col-span-2">
            <div className="text-xs opacity-80">Recent price trend</div>
            <div className="mt-1">
              <CompanySparkline series={series} />
            </div>
          </div>
          <div className="glass-card p-3">
            <div className="text-xs opacity-80">Snapshot</div>
            <div className="mt-1 text-sm" style={{ color: "#0B2545" }}>
              <div>Price: {price?.price?.close_price ?? "—"}</div>
              <div>ESG: {esg?.overall_score?.toFixed ? esg.overall_score.toFixed(2) : "—"}</div>
              <div>P/E: {ind?.indicators?.pe_ratio ?? "—"}</div>
              <div>ROE: {ind?.indicators?.return_on_equity ?? "—"}</div>
            </div>
          </div>
        </div>
        <div className="px-4 pb-4 grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div className="glass-card p-3">
            <div className="text-[11px]" style={{ color: "#374151" }}>Sector</div>
            <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{company?.sector || "—"}</div>
          </div>
          <div className="glass-card p-3">
            <div className="text-[11px]" style={{ color: "#374151" }}>Industry</div>
            <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{company?.industry || "—"}</div>
          </div>
          <div className="glass-card p-3">
            <div className="text-[11px]" style={{ color: "#374151" }}>Market Cap</div>
            <div className="text-sm font-semibold" style={{ color: "#0B2545" }}>{formatCompact(company?.market_cap)}</div>
          </div>
        </div>
        {summary && (
          <div className="px-4 pb-4">
            <div className="glass-card p-3">
              <div className="text-xs opacity-80">Daily summary</div>
              <div className="grid grid-cols-3 gap-2 mt-2 text-sm" style={{ color: "#0B2545" }}>
                <div>Δ: {(summary.summary.price_change >= 0 ? "+" : "") + summary.summary.price_change.toFixed(2)}</div>
                <div>Δ%: {(summary.summary.price_change_percent >= 0 ? "+" : "") + summary.summary.price_change_percent.toFixed(2)}%</div>
                <div>Volume: {summary.summary.volume.toLocaleString()}</div>
              </div>
            </div>
          </div>
        )}
        {loading && (
          <div className="px-4 pb-4">
            <div className="hero-progress"><span /></div>
          </div>
        )}
      </div>
    </div>
  );
}

function formatCompact(n?: number) {
  if (!n && n !== 0) return "—";
  if (n >= 1e12) return (n / 1e12).toFixed(2) + 'T';
  if (n >= 1e9) return (n / 1e9).toFixed(2) + 'B';
  if (n >= 1e6) return (n / 1e6).toFixed(2) + 'M';
  if (n >= 1e3) return (n / 1e3).toFixed(2) + 'K';
  return n.toFixed(0);
}


