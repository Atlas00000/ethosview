import React from "react";
import type { FinancialIndicatorsResponse, LatestPriceResponse } from "../../types/api";

export function FinancialSnapshot({ ind, price }: { ind: FinancialIndicatorsResponse; price: LatestPriceResponse }) {
  const i = ind?.indicators || {};
  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <h2 className="text-xl font-semibold mb-3 text-gradient">Company snapshot</h2>
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
        <div className="glass-card p-3">
          <div className="text-xs" style={{ color: "#374151" }}>P/E</div>
          <div className="text-lg font-medium" style={{ color: "#0B2545" }}>{i.pe_ratio ?? "—"}</div>
        </div>
        <div className="glass-card p-3">
          <div className="text-xs" style={{ color: "#374151" }}>ROE</div>
          <div className="text-lg font-medium" style={{ color: "#0B2545" }}>{i.return_on_equity ?? "—"}</div>
        </div>
        <div className="glass-card p-3">
          <div className="text-xs" style={{ color: "#374151" }}>Margin</div>
          <div className="text-lg font-medium" style={{ color: "#0B2545" }}>{i.profit_margin ?? "—"}</div>
        </div>
        <div className="glass-card p-3">
          <div className="text-xs" style={{ color: "#374151" }}>Price</div>
          <div className="text-lg font-medium" style={{ color: "#0B2545" }}>{price?.price?.close_price ?? "—"}</div>
        </div>
      </div>
    </section>
  );
}


