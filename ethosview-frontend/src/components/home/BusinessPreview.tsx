"use client";
import React from "react";
import type { DashboardResponse } from "../../types/api";

export function BusinessPreview({ dashboard }: { dashboard: DashboardResponse }) {
  const s = dashboard.summary;
  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <h2 className="text-xl font-semibold mb-2">Business dashboard preview</h2>
      <div className="rounded-lg border p-4 grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div>
          <div className="text-sm text-gray-500">Total companies</div>
          <div className="text-2xl font-medium">{s.total_companies.toLocaleString()}</div>
        </div>
        <div>
          <div className="text-sm text-gray-500">Sectors</div>
          <div className="text-2xl font-medium">{s.total_sectors}</div>
        </div>
        <div>
          <div className="text-sm text-gray-500">Avg ESG</div>
          <div className="text-2xl font-medium">{s.avg_esg_score.toFixed(2)}</div>
        </div>
      </div>
    </section>
  );
}


