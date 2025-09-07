"use client";
import React, { useEffect, useState } from "react";
import { api } from "../../services/api";

const nav = [
  { href: "#hero", label: "Home" },
  { href: "#market", label: "Market" },
  { href: "#esg", label: "ESG" },
  { href: "#featured", label: "Featured" },
  { href: "#sectors", label: "Sectors" },
  
];

export function Header() {
  const [uptime, setUptime] = useState<string>("");
  const [alerts, setAlerts] = useState<number>(0);
  useEffect(() => {
    let alive = true;
    // Fetch a couple of lightweight backend signals: health/live and alerts count
    Promise.all([
      fetch(`${api.API_BASE_URL ?? ''}/health/live`).then((r) => r.ok ? r.json() : null).catch(() => null),
      api.alerts(false).catch(() => ({ alerts: [], count: 0 } as any)),
    ]).then(([live, a]) => {
      if (!alive) return;
      if (live?.status) setUptime("Live");
      setAlerts((a as any)?.count ?? 0);
    });
    return () => { alive = false; };
  }, []);

  return (
    <header className="sticky top-0 z-40 backdrop-blur border-b" style={{ background: "rgba(255,255,255,0.7)", borderColor: "rgba(0,0,0,0.06)" }}>
      <div className="max-w-6xl mx-auto px-4 h-14 flex items-center justify-between">
        <a href="#hero" className="font-semibold text-gradient hover-lift">EthosView</a>
        <nav className="flex items-center gap-5 text-sm">
          {nav.map((n) => (
            <a key={n.href} href={n.href} className="text-gray-700 hover:text-black transition-colors tilt-hover" style={{ color: "#374151" }}>
              {n.label}
            </a>
          ))}
        </nav>
        <div className="hidden sm:flex items-center gap-3 text-xs" style={{ color: "#374151" }}>
          {uptime && <span className="px-2 py-0.5 rounded-full bg-white/60">API: {uptime}</span>}
          <span className="px-2 py-0.5 rounded-full bg-white/60">Alerts: {alerts}</span>
        </div>
      </div>
    </header>
  );
}


