"use client";
import React from "react";
import type { MarketLatestResponse } from "../../types/api";

export function MarketBar({ market }: { market: MarketLatestResponse }) {
  const m = market.market_data;
  return (
    <section className="animated-gradient">
      <div className="max-w-6xl mx-auto px-4 py-3 flex flex-wrap gap-6 items-center text-sm text-white" style={{ transition: "filter var(--dur-hover) var(--ease-enter)" }}>
        <div className="font-medium">Market Snapshot</div>
        <div>S&amp;P 500: {m.sp500_close ?? ""}</div>
        <div>NASDAQ: {m.nasdaq_close ?? ""}</div>
        <div>DOW: {m.dow_close ?? ""}</div>
        <div>VIX: {m.vix_close ?? ""}</div>
        <div>10Y: {m.treasury_10y ?? ""}</div>
      </div>
    </section>
  );
}


