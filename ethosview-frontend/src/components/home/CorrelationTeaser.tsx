import React from "react";
import type { CorrelationResponse } from "../../types/api";

function formatValue(v: unknown): string {
  const num = typeof v === "number" && Number.isFinite(v) ? v : NaN;
  return Number.isFinite(num) ? num.toFixed(2) : "—";
}

export function CorrelationTeaser({ corr }: { corr: CorrelationResponse | Record<string, any> }) {
  const items = [
    { label: "ESG vs Market Cap", value: (corr as any)?.esg_market_cap_corr },
    { label: "ESG vs P/E", value: (corr as any)?.esg_pe_corr },
    { label: "ESG vs ROE", value: (corr as any)?.esg_roe_corr },
    { label: "ESG vs Profit", value: (corr as any)?.esg_profit_corr },
  ];
  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <h2 className="text-xl font-semibold mb-3 text-gradient">ESG ↔ Financial correlation</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-3">
        {items.map((it) => (
          <div key={it.label} className="glass-card p-3">
            <div className="text-xs" style={{ color: "#374151" }}>{it.label}</div>
            <div className="text-lg font-medium" style={{ color: "#0B2545" }}>{formatValue(it.value)}</div>
          </div>
        ))}
      </div>
    </section>
  );
}


