"use client";
import React from "react";
import type { PerformanceMetric, LatestPriceResponse, LatestESGResponse } from "../../types/api";

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

  return (
    <section className="max-w-6xl mx-auto px-4 py-12 relative overflow-hidden">
      <div className="blob" style={{ bottom: -80, left: -60, width: 200, height: 200, background: "rgba(30,106,225,0.12)", borderRadius: 9999 }} />
      <h2 className="text-xl font-semibold mb-4" style={{ color: "#0B2545" }}>Featured companies</h2>
      <div className="flex gap-4 overflow-x-auto pb-2 edge-fade snap-x snap-mandatory">
        {items.map((it) => (
          <div key={it.company_id} className="min-w-[240px] shrink-0 glass-card p-4 hover-lift snap-start">
            <div className="text-sm" style={{ color: "#6B7280" }}>#{it.company_id}</div>
            <div className="text-lg font-medium" style={{ color: "#0B2545" }}>{it.company_name}</div>
            {typeof it.price === "number" ? (
              <div className="text-sm mt-1" style={{ color: "#0B2545" }}>Price: {it.price.toFixed(2)}</div>
            ) : (
              <div className="text-sm mt-1" style={{ color: "#6B7280" }}>Price loading…</div>
            )}
            {typeof it.esg === "number" ? (
              <div className="text-sm" style={{ color: "#1D9A6C" }}>ESG: {it.esg.toFixed(2)}</div>
            ) : (
              <div className="text-sm" style={{ color: "#6B7280" }}>ESG loading…</div>
            )}
          </div>
        ))}
      </div>
    </section>
  );
}


