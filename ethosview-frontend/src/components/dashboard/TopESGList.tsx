"use client";
import React from "react";
import type { PerformanceMetric } from "../../types/api";

export function TopESGList({ items }: { items: PerformanceMetric[] }) {
  return (
    <div className="divide-y">
      {items.map((it) => (
        <div key={it.company_id} className="py-2 flex items-center justify-between">
          <div className="truncate pr-3">
            <div className="text-sm font-medium truncate" style={{ color: "#0B2545" }}>{it.company_name}</div>
            <div className="text-xs" style={{ color: "#374151" }}>Rank #{it.rank} • {it.percentile.toFixed(1)}%</div>
          </div>
          <div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{it.value.toFixed(2)}</div>
        </div>
      ))}
    </div>
  );
}


