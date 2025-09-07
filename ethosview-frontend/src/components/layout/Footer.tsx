"use client";
import React, { useEffect, useState } from "react";
import { api, API_BASE_URL } from "../../services/api";

export function Footer() {
  const [metrics, setMetrics] = useState<{ timestamp?: string } | null>(null);
  useEffect(() => {
    let alive = true;
    fetch(`${API_BASE_URL}/metrics`).then(r => (r.ok ? r.json() : null)).then((m) => { if (alive) setMetrics(m); }).catch(() => {});
    return () => { alive = false; };
  }, []);
  return (
    <footer className="mt-16 border-t" style={{ borderColor: "rgba(0,0,0,0.06)" }}>
      <div className="max-w-6xl mx-auto px-4 py-8 text-sm flex flex-col sm:flex-row items-center justify-between gap-3">
        <div className="text-gradient font-semibold">EthosView</div>
        <div className="text-gray-600">ESG and Financial Analytics</div>
        <div className="flex items-center gap-3 text-gray-500">
          <span>© {new Date().getFullYear()} EthosView</span>
          {metrics && <span className="hidden sm:inline">• Last metrics: {new Date(metrics.timestamp || Date.now()).toLocaleString()}</span>}
        </div>
      </div>
    </footer>
  );
}


