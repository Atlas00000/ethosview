"use client";
import React, { useEffect, useState } from "react";
import { api, API_BASE_URL } from "../../services/api";

export function QuickWidget() {
  const [up, setUp] = useState<boolean | null>(null);
  const [alerts, setAlerts] = useState<number>(0);
  useEffect(() => {
    let alive = true;
    Promise.all([
      fetch(`${API_BASE_URL}/health/live`).then(r => r.ok ? r.json() : null).catch(() => null),
      api.alerts(false).catch(() => ({ count: 0 } as any))
    ]).then(([live, a]) => {
      if (!alive) return;
      setUp(!!live?.status);
      setAlerts((a as any)?.count ?? 0);
    });
    return () => { alive = false; };
  }, []);

  return (
    <div className="fixed bottom-5 right-5 z-50 flex flex-col items-end gap-2">
      <a href="#hero" aria-label="Back to top" className="glass-card glow-ping hover-lift px-3 py-2 rounded-full text-sm" style={{ color: "#0B2545" }}>
        ↑ Top
      </a>
      <div className="glass-card px-3 py-2 rounded-full text-xs flex items-center gap-2" style={{ color: "#374151" }}>
        <span style={{ width: 8, height: 8, borderRadius: 9999, background: up ? '#1D9A6C' : '#E25555', display: 'inline-block' }} />
        <span>{up ? 'Live' : 'Offline'}</span>
        <span>• Alerts {alerts}</span>
      </div>
    </div>
  );
}


