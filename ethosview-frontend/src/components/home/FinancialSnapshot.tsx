"use client";
import React from "react";
import type { FinancialIndicatorsResponse, LatestPriceResponse, CompanyFinancialSummaryResponse, CompanyPricesResponse } from "../../types/api";
import { CompanySparkline } from "./MarketSparkline";

export function FinancialSnapshot({ ind, price, summary, series }: { ind: FinancialIndicatorsResponse; price: LatestPriceResponse; summary?: CompanyFinancialSummaryResponse | null; series?: CompanyPricesResponse | null }) {
  const i = ind?.indicators || {};
  const s = summary?.summary;
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
      {(series || s) && (
        <div className="mt-3 glass-card p-3 animate-fade-in-up">
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 items-center">
            <div className="sm:col-span-2">
              <div className="text-xs opacity-80">Recent price trend</div>
              <div className="mt-1">
                <CompanySparkline series={(series as any) ?? null} />
              </div>
            </div>
            <div>
              <div className="text-xs opacity-80">Summary</div>
              <div className="mt-1 text-sm" style={{ color: "#0B2545" }}>
                {s ? (
                  <>
                    <div>Change: {(s.price_change >= 0 ? '+' : '') + s.price_change.toFixed(2)} ({(s.price_change_percent >= 0 ? '+' : '') + s.price_change_percent.toFixed(2)}%)</div>
                    <div>Volume: {s.volume.toLocaleString()}</div>
                    <div>Date: {s.date}</div>
                  </>
                ) : (
                  <div>—</div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}
    </section>
  );
}


